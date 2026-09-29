package web

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func repoRegistryPath() string {
	return qweDir() + "/repositories"
}

func gitRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func loadRepos() []string {
	data, err := os.ReadFile(repoRegistryPath())
	if err != nil {
		return nil
	}
	var repos []string
	for _, line := range strings.Split(string(data), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			repos = append(repos, line)
		}
	}
	return repos
}

func saveRepos(repos []string) {
	os.WriteFile(repoRegistryPath(), []byte(strings.Join(repos, "\n")+"\n"), 0600)
}

// KnownRepos returns the registry of repositories gocommit knows about.
func KnownRepos() []string {
	return loadRepos()
}

// DiscoverSiblingRepos returns known repos plus any git repositories found
// alongside them: for each known repo's parent directory, every subdirectory
// containing a .git is included. This picks up e.g. all repos under a shared
// "sg" workspace folder once any one of them is known.
func DiscoverSiblingRepos(known []string) []string {
	seen := make(map[string]bool, len(known))
	repos := make([]string, 0, len(known))
	for _, r := range known {
		if seen[r] {
			continue
		}
		if _, err := os.Stat(filepath.Join(r, ".git")); err != nil {
			continue
		}
		seen[r] = true
		repos = append(repos, r)
	}
	visitedDirs := make(map[string]bool)
	for _, r := range known {
		dir := filepath.Dir(r)
		if visitedDirs[dir] {
			continue
		}
		visitedDirs[dir] = true
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			candidate := filepath.Join(dir, e.Name())
			if seen[candidate] {
				continue
			}
			if _, err := os.Stat(filepath.Join(candidate, ".git")); err != nil {
				continue
			}
			seen[candidate] = true
			repos = append(repos, candidate)
		}
	}
	sort.Strings(repos)
	return repos
}

func RegisterCurrentRepo() {
	root, err := gitRoot()
	if err != nil {
		return
	}
	repos := loadRepos()
	for _, r := range repos {
		if r == root {
			return
		}
	}
	repos = append(repos, root)
	saveRepos(repos)
}

func RepositoriesController(w http.ResponseWriter, r *http.Request) {
	repos := DiscoverSiblingRepos(loadRepos())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func CurrentRepoController(w http.ResponseWriter, r *http.Request) {
	root, _ := gitRoot()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(root)
}

func SwitchRepoController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Path == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := os.Chdir(body.Path); err != nil {
		http.Error(w, "failed to switch repo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"path": body.Path})
}