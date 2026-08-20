package commit

import (
	"os/exec"
	"regexp"
	"strings"
)

func GetRepoInfo() (map[string]interface{}, error) {
	return map[string]interface{}{"message": "niente da vedere qui"}, nil
}

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

// GetUncommittedFiles returns the list of uncommitted files in porcelain format
func GetUncommittedFiles() ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var files []string
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// IsWorkingDirectoryDirty returns true if there are uncommitted changes
func IsWorkingDirectoryDirty() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(output)) != "", nil
}

func CheckoutBranch(branch string) error {
	dirty, err := IsWorkingDirectoryDirty()
	if err != nil {
		return err
	}
	if dirty {
		return &CheckoutError{"Cannot switch branch: working directory is dirty. Please commit or stash your changes first."}
	}
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

type CheckoutError struct {
	msg string
}

func (e *CheckoutError) Error() string {
	return e.msg
}

// RecentFeatureNames parses git log subjects and returns unique scopes from
// conventional commits (e.g. "feat(scope): msg"), most-recent first, up to n.
func RecentFeatureNames(n int) []string {
	out, err := exec.Command("git", "log", "--pretty=format:%s", "-100").Output()
	if err != nil {
		return nil
	}
	re := regexp.MustCompile(`^\w+\(([^)]+)\):`)
	seen := map[string]bool{}
	var names []string
	for _, line := range strings.Split(string(out), "\n") {
		m := re.FindStringSubmatch(strings.TrimSpace(line))
		if len(m) < 2 {
			continue
		}
		scope := m[1]
		if !seen[scope] {
			seen[scope] = true
			names = append(names, scope)
			if len(names) == n {
				break
			}
		}
	}
	return names
}

// GetFileDiff returns the diff for a single file vs HEAD.
// For untracked files (status "??") it shows the full content as an addition.
func GetFileDiff(file string) (string, error) {
	out, err := exec.Command("git", "diff", "HEAD", "--", file).Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return string(out), nil
	}
	// Fallback for new untracked files
	out, err = exec.Command("git", "diff", "--no-index", "/dev/null", file).Output()
	if err != nil && len(out) == 0 {
		return "", err
	}
	return string(out), nil
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

func BranchHasUpstream() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	err := cmd.Run()
	if err != nil {
		return false, nil
	}
	return true, nil
}

func RunGitPush() error {
	cmd := exec.Command("git", "push")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func RunGitPushSetUpstream(branch string) error {
	cmd := exec.Command("git", "push", "--set-upstream", "origin", branch)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// RunGitMerge merges the given branch into the current branch with --no-ff.
func RunGitMerge(branch string) error {
	cmd := exec.Command("git", "merge", "--no-ff", branch)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// DeleteLocalBranch deletes a local branch that has been fully merged.
func DeleteLocalBranch(branch string) error {
	cmd := exec.Command("git", "branch", "-d", branch)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
