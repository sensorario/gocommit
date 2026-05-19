package web

import (
	"encoding/json"
	"gocommit/commit"
	"net/http"
	"os"
)

func BranchesController(w http.ResponseWriter, r *http.Request) {
	branches, err := commit.ListLocalBranches()
	if err != nil {
		http.Error(w, "failed to list branches", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branches)
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

func MainPageController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html, err := os.ReadFile("internal/web/main.html")
	if err != nil {
		http.Error(w, "Could not load HTML", http.StatusInternalServerError)
		return
	}
	w.Write(html)
}
