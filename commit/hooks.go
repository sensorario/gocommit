package commit

import (
	"fmt"
	"os/exec"
)

func RunHook(hook string) error {
	if hook == "" {
		return nil
	}
	fmt.Printf("Running %s: %s\n", hook, hook)
	shellCmd := exec.Command("sh", "-c", hook)
	shellCmd.Stdout = nil
	shellCmd.Stderr = nil
	return shellCmd.Run()
}
