package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type PullRequest struct {
	Title  string `json:"title"`
	Number int    `json:"number"`
	URL    string `json:"url"`
}

func PullRequestsController(w http.ResponseWriter, r *http.Request) {
	remotes, err := ListRemotes()
	if err != nil || len(remotes) == 0 {
		http.Error(w, "no remotes found", http.StatusInternalServerError)
		return
	}

	remoteURL := remotes[0].URL
	for _, remote := range remotes {
		if remote.Name == "origin" {
			remoteURL = remote.URL
			break
		}
	}

	prs, err := fetchPullRequests(remoteURL)
	if err != nil {
		http.Error(w, "failed to fetch pull requests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prs)
}

func fetchPullRequests(remoteURL string) ([]PullRequest, error) {
	owner, repo, provider := parseRemoteURL(remoteURL)
	if provider == "" {
		return nil, fmt.Errorf("unsupported remote: %s", remoteURL)
	}
	switch provider {
	case "github":
		return fetchGitHubPRs(owner, repo)
	case "bitbucket":
		return fetchBitbucketPRs(owner, repo)
	}
	return nil, fmt.Errorf("unknown provider")
}

func parseRemoteURL(url string) (owner, repo, provider string) {
	type pattern struct {
		prefix   string
		provider string
	}
	patterns := []pattern{
		{"git@github.com:", "github"},
		{"https://github.com/", "github"},
		{"git@bitbucket.org:", "bitbucket"},
		{"https://bitbucket.org/", "bitbucket"},
	}
	url = strings.TrimSpace(url)
	for _, p := range patterns {
		if strings.HasPrefix(url, p.prefix) {
			rest := strings.TrimPrefix(url, p.prefix)
			rest = strings.TrimSuffix(rest, ".git")
			parts := strings.SplitN(rest, "/", 2)
			if len(parts) == 2 {
				return parts[0], parts[1], p.provider
			}
		}
	}
	return "", "", ""
}

func fetchGitHubPRs(owner, repo string) ([]PullRequest, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=open", owner, repo)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var raw []struct {
		Title   string `json:"title"`
		Number  int    `json:"number"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	prs := make([]PullRequest, len(raw))
	for i, p := range raw {
		prs[i] = PullRequest{Title: p.Title, Number: p.Number, URL: p.HTMLURL}
	}
	return prs, nil
}

func fetchBitbucketPRs(workspace, repo string) ([]PullRequest, error) {
	apiURL := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests?state=OPEN", workspace, repo)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	if token := os.Getenv("BITBUCKET_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else if user := os.Getenv("BITBUCKET_USERNAME"); user != "" {
		if pass := os.Getenv("BITBUCKET_APP_PASSWORD"); pass != "" {
			req.SetBasicAuth(user, pass)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Bitbucket API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var raw struct {
		Values []struct {
			Title string `json:"title"`
			ID    int    `json:"id"`
			Links struct {
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
			} `json:"links"`
		} `json:"values"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	prs := make([]PullRequest, len(raw.Values))
	for i, p := range raw.Values {
		prs[i] = PullRequest{Title: p.Title, Number: p.ID, URL: p.Links.HTML.Href}
	}
	return prs, nil
}