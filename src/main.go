package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"

	"commit"
)

var Version = "dev"

// ...existing code...

func main() {

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-v" || arg == "--version" {
			fmt.Println(Version)
			return
		}
		commit.PrintRed("\u001f Non sono accettati parametri ad eccezione di -v o --version")
		os.Exit(1)
	}

	configPath := "gocommit.conf.json"
	if err := commit.EnsureConfig(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "[ERR1001] Error creating config: %v\n", err)
		os.Exit(1)
	}

	if err := commit.RunGitAdd(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERR1002] Error during git add: %v\n", err)
		os.Exit(1)
	}

	config, err := commit.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERR1003] Error loading config: %v\n", err)
		os.Exit(1)
	}
	if err := commit.RunHook(config.OnBeforeCommit); err != nil {
		fmt.Fprintf(os.Stderr, "[ERR1004] Error running onBeforeCommit: %v\n", err)
		os.Exit(1)
	}

	prompt := promptui.Select{
		Label: "Commit type",
		Items: []string{"feat", "fix", "chore", "refactor"},
	}
	_, commitType, err := prompt.Run()
	if err != nil {
		fmt.Printf("[ERR1005] Prompt failed: %v\n", err)
		return
	}

	branchName, err := commit.CurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERR1006] Error getting branch: %v\n", err)
		os.Exit(1)
	}
	ticket := commit.ExtractJiraTicket(branchName)
	feature := commit.AskFeatureName(ticket)
	message := commit.AskCommitMessage()

	if err := commit.RunHook(config.OnAfterCommit); err != nil {
		fmt.Fprintf(os.Stderr, "[ERR1007] Error running onAfterCommit: %v\n", err)
		os.Exit(1)
	}

	fullMessage := commitType + "(" + feature + "): " + message
	fmt.Println("Running: git commit -m \"" + fullMessage + "\"")
	if err := commit.RunGitCommit(fullMessage); err != nil {
		fmt.Fprintf(os.Stderr, "[ERR1008] Error during git commit: %v\n", err)
		os.Exit(1)
	}

	exists, err := commit.RemoteExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERR1009] Failed to check remote: %v\n", err)
		os.Exit(1)
	}

	if !exists {
		commit.PrintRed("[ERR1010] \u001f Nessun remote presente. Usa 'git remote add origin <url>' per aggiungerne uno.")
		os.Exit(1)
	}

	fmt.Println("Running: git push")
	if err := commit.RunGitPush(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERR1011] Error during git push: %v\n", err)
		os.Exit(1)
	}
}
