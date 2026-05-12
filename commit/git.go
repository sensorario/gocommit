package commit

import (
	"os/exec"
	"strings"
)

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
