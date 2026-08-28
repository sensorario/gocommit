package branch

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"

	"gocommit/commit"
)

func printUncommittedFiles(files []string) {
	const width = 42
	title := " Uncommitted files "
	top := "┌─" + title + strings.Repeat("─", width-len(title)-2) + "┐"
	bottom := "└" + strings.Repeat("─", width) + "┘"
	fmt.Println(top)
	for _, f := range files {
		line := "│  " + f
		padding := width - len([]rune(f)) - 3
		if padding < 0 {
			padding = 0
		}
		fmt.Println(line + strings.Repeat(" ", padding) + "│")
	}
	fmt.Println(bottom)
}

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
	if files, err := commit.GetUncommittedFiles(); err == nil && len(files) > 0 {
		printUncommittedFiles(files)
	}
	current, _ := commit.CurrentBranch()
	items := make([]string, len(branches))
	for i, b := range branches {
		if b == current {
			items[i] = b + "  (current)"
		} else {
			items[i] = b
		}
	}
	prompt := promptui.Select{
		Label: "Select branch",
		Items: items,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		os.Exit(1)
	}
	selectedBranch := branches[idx]
	   if err := commit.CheckoutBranch(selectedBranch); err != nil {
		   fmt.Fprintf(os.Stderr, "Error switching branch: %v\n", err)
		   if _, ok := err.(*commit.CheckoutError); ok {
			   os.Exit(2)
		   }
		   os.Exit(1)
	   }
	   commit.PrintGreen("Switched to branch: " + selectedBranch)
}
