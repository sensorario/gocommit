package web

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func repoRegistryPath() string {
	home, _ := os.UserHomeDir()
	return home + "/.sensorario-qwe/repositories"
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
	repos := loadRepos()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func CurrentRepoController(w http.ResponseWriter, r *http.Request) {
	root, err := gitRoot()
	if err != nil {
		http.Error(w, "not a git repository", http.StatusInternalServerError)
		return
	}
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