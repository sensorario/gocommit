package web

import (
	"bufio"
	"os"
	"regexp"
)

type Remote struct {
	Name string
	URL  string
}

// ListRemotes parses .git/config and returns all remotes with their URLs
func ListRemotes() ([]Remote, error) {
	file, err := os.Open(".git/config")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	remotes := []Remote{}
	var current string
	remoteRe := regexp.MustCompile(`^\[remote \"(.+)\"\]`)
	urlRe := regexp.MustCompile(`^\s*url\s*=\s*(.+)$`)
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if m := remoteRe.FindStringSubmatch(line); m != nil {
			current = m[1]
			continue
		}
		if current != "" {
			if m := urlRe.FindStringSubmatch(line); m != nil {
				remotes = append(remotes, Remote{Name: current, URL: m[1]})
				current = ""
			}
		}
	}
	return remotes, scanner.Err()
}
