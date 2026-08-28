package web

import (
	"encoding/json"
	"net/http"

	"gocommit/commit"
)

type RepoVersion struct {
	Path    string
	Tech    string
	Version string
	Managed bool
}

func RepoVersionsController(w http.ResponseWriter, r *http.Request) {
	repos := DiscoverSiblingRepos(loadRepos())
	versions := make([]RepoVersion, 0, len(repos))
	for _, repo := range repos {
		tech := commit.DetectProjectTech(repo)
		version := commit.ProjectVersion(repo, tech)
		versions = append(versions, RepoVersion{
			Path:    repo,
			Tech:    string(tech),
			Version: version,
			Managed: commit.VersionManaged(version),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)
}
