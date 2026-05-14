package commit

import (
	"os/exec"
	"strings"
)

// ListLocalBranches returns a slice of all local git branch names
func ListLocalBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	var branches []string
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

func CurrentBranch() (string, error) {
	branchNameBytes, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(branchNameBytes)), nil
}

func RemoteExists() (bool, error) {
	cmd := exec.Command("git", "remote")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(output)) != "", nil
}

func ListBranches() ([]string, error) {
	output, err := exec.Command("git", "branch", "--format=%(refname:short)").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var branches []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			branches = append(branches, l)
		}
	}
	return branches, nil
}

func CheckoutBranch(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func RunGitAdd() error {
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func RunGitCommit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func RunGitPush() error {
	cmd := exec.Command("git", "push")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
