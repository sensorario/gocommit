package srcutils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GitRemoteExists() (bool, error) {
	cmd := exec.Command("git", "remote")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(output)) != "", nil
}

func PrintRed(msg string) {
	fmt.Printf("\033[31m%s\033[0m\n", msg)
}
