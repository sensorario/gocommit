package check

import (
	"fmt"
	"os"

	"gocommit/commit"
)

func Run() {
	configPath := "gocommit.conf.json"
	added, err := commit.CheckConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking config: %v\n", err)
		os.Exit(1)
	}
	if len(added) == 0 {
		commit.PrintGreen("Config OK: all variables are present in " + configPath)
	} else {
		for _, key := range added {
			commit.PrintYellow("Added missing variable with default: " + key)
		}
	}
}
