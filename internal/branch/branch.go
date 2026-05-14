package branch

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"

	"gocommit/commit"
)

func Run() {
	branches, err := commit.ListLocalBranches()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing branches: %v\n", err)
		os.Exit(1)
	}
	if len(branches) == 0 {
		fmt.Fprintf(os.Stderr, "No local branches found.\n")
		os.Exit(1)
	}
	prompt := promptui.Select{
		Label: "Select branch",
		Items: branches,
	}
	_, selectedBranch, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		os.Exit(1)
	}
	if err := commit.CheckoutBranch(selectedBranch); err != nil {
		fmt.Fprintf(os.Stderr, "Error switching branch: %v\n", err)
		os.Exit(1)
	}
	commit.PrintGreen("Switched to branch: " + selectedBranch)
}
