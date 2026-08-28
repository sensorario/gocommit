package versions

import (
	"fmt"
	"path/filepath"

	"gocommit/commit"
	"gocommit/internal/web"
)

func Run() {
	repos := web.DiscoverSiblingRepos(web.KnownRepos())
	if len(repos) == 0 {
		fmt.Println("No known repositories.")
		return
	}
	for _, repo := range repos {
		name := filepath.Base(repo)
		tech := commit.DetectProjectTech(repo)
		version := commit.ProjectVersion(repo, tech)

		techLabel := string(tech)
		if techLabel == "" {
			techLabel = "-"
		}
		versionLabel := version
		if versionLabel == "" {
			versionLabel = "n/a"
		}
		managedLabel := "not managed"
		if commit.VersionManaged(version) {
			managedLabel = "managed"
		}
		fmt.Printf("%-30s %-6s %-10s %s\n", name, techLabel, versionLabel, managedLabel)
	}
}
