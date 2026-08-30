package web

import (
	"encoding/json"
	"gocommit/commit"
	"net/http"
	"os"
)

var DevMode = false

func init() {
       if os.Getenv("DEV_MODE") == "1" {
	       DevMode = true
       }
}

// Restituisce la versione corrente
func VersionController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(GetVersion()))
}

// Restituisce tutti i remoti come JSON
func RemotesController(w http.ResponseWriter, r *http.Request) {
	remotes, err := ListRemotes()
	if err != nil {
		http.Error(w, "Could not read remotes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(remotes)
}

func BranchesController(w http.ResponseWriter, r *http.Request) {
	branches, err := commit.ListLocalBranches()
	if err != nil {
		http.Error(w, "failed to list branches", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branches)
}

func CurrentBranchController(w http.ResponseWriter, r *http.Request) {
	branch, err := commit.CurrentBranch()
	if err != nil {
		http.Error(w, "failed to get current branch", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branch)
}

func InfoController(w http.ResponseWriter, r *http.Request) {
	info, err := commit.GetRepoInfo()
	if err != nil {
		http.Error(w, "failed to get repo info", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func CheckoutController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Branch string `json:"branch"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Branch == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := commit.CheckoutBranch(body.Branch); err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(*commit.CheckoutError); ok {
			status = http.StatusConflict
		}
		http.Error(w, "checkout failed: "+err.Error(), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"branch": body.Branch})
}

func MergeIntoNextController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Branch string `json:"branch"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Branch == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if body.Branch == "next" {
		http.Error(w, "cannot merge next into itself", http.StatusBadRequest)
		return
	}
	if err := commit.CheckoutBranch("next"); err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(*commit.CheckoutError); ok {
			status = http.StatusConflict
		}
		http.Error(w, "checkout failed: "+err.Error(), status)
		return
	}
	if err := commit.RunGitMerge(body.Branch); err != nil {
		http.Error(w, "merge failed: "+err.Error(), http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"branch": "next", "merged": body.Branch})
}

func MergedBranchesController(w http.ResponseWriter, r *http.Request) {
	branches, err := commit.MergedLocalBranches("next")
	if err != nil {
		http.Error(w, "failed to list merged branches", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branches)
}

func DeleteBranchController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Branch string `json:"branch"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Branch == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if body.Branch == "next" {
		http.Error(w, "cannot delete next", http.StatusBadRequest)
		return
	}
	if err := commit.DeleteLocalBranch(body.Branch); err != nil {
		http.Error(w, "delete failed: "+err.Error(), http.StatusConflict)
		return
	}
	remoteExists, err := commit.RemoteBranchExists(body.Branch)
	if err == nil && remoteExists {
		commit.DeleteRemoteBranch(body.Branch)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"deleted": body.Branch})
}

func ShutdownController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "shutting down"})
	shutdownCh <- struct{}{}
}

func ModifiedFilesController(w http.ResponseWriter, r *http.Request) {
	files, err := commit.GetUncommittedFiles()
	if err != nil {
		http.Error(w, "failed to get git status", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func DiffController(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "missing file param", http.StatusBadRequest)
		return
	}
	diff, err := commit.GetFileDiff(file)
	if err != nil {
		http.Error(w, "failed to get diff", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(diff))
}

func MainPageController(w http.ResponseWriter, r *http.Request) {
	       w.Header().Set("Content-Type", "text/html; charset=utf-8")
	       var html []byte
	       var err error
	       if DevMode {
		       html, err = os.ReadFile("internal/web/main.html")
		       if err != nil {
			       http.Error(w, "Could not load HTML", http.StatusInternalServerError)
			       return
		       }
	       } else {
		       html = EmbeddedHTML
	       }
	       w.Write(html)
}
