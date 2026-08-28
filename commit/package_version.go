package commit

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func readJSONVersion(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", err
	}
	return pkg.Version, nil
}

// PackageVersion reads the "version" field from package.json in repoPath.
func PackageVersion(repoPath string) (string, error) {
	return readJSONVersion(filepath.Join(repoPath, "package.json"))
}

// ComposerVersion reads the "version" field from composer.json in repoPath.
func ComposerVersion(repoPath string) (string, error) {
	return readJSONVersion(filepath.Join(repoPath, "composer.json"))
}
