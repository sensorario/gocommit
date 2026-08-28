package commit

import (
	"os"
	"path/filepath"
	"strings"
)

type ProjectTech string

const (
	TechNode    ProjectTech = "node"
	TechGo      ProjectTech = "go"
	TechPHP     ProjectTech = "php"
	TechUnknown ProjectTech = ""
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DetectProjectTech inspects repoPath for well-known marker files and
// reports which technology the project appears to use.
func DetectProjectTech(repoPath string) ProjectTech {
	if fileExists(filepath.Join(repoPath, "package.json")) {
		return TechNode
	}
	if fileExists(filepath.Join(repoPath, "go.mod")) {
		return TechGo
	}
	if fileExists(filepath.Join(repoPath, "composer.json")) {
		return TechPHP
	}
	return TechUnknown
}

// ProjectVersion reads repoPath's version using the convention for tech:
// the "version" field of package.json for Node, the VERSION file for Go
// (this project's own convention), and the "version" field of
// composer.json for PHP.
func ProjectVersion(repoPath string, tech ProjectTech) string {
	switch tech {
	case TechNode:
		version, _ := PackageVersion(repoPath)
		return version
	case TechGo:
		data, err := os.ReadFile(filepath.Join(repoPath, "VERSION"))
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(data))
	case TechPHP:
		version, _ := ComposerVersion(repoPath)
		return version
	default:
		return ""
	}
}

// VersionManaged reports whether version looks like a real, maintained
// version rather than an unset or default placeholder.
func VersionManaged(version string) bool {
	return version != "" && version != "0.0.0"
}
