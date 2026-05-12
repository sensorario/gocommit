package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"

	"gocommit/commit"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-v" || arg == "--version" {
			fmt.Println(Version)
			return
		}
		if arg == "help" || arg == "--help" || arg == "-h" {
			fmt.Println("Usage: qwe [command]")
			fmt.Println()
			fmt.Println("Commands:")
			fmt.Printf("  %-20s %s\n", "(no args)", "Start the interactive commit wizard")
			fmt.Printf("  %-20s %s\n", "check", "Verify gocommit.conf.json exists and all variables are present; adds missing ones with defaults")
			fmt.Printf("  %-20s %s\n", "help, -h, --help", "Show this help message")
			fmt.Printf("  %-20s %s\n", "-v, --version", "Print the current version")
			return
		}
		if arg == "check" {
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
			return
		}
		commit.PrintRed("Unknown command: " + arg + ". Run 'qwe help' to see available commands.")
		os.Exit(1)
	}

	configPath := "gocommit.conf.json"
	if err := commit.EnsureConfig(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config: %v\n", err)
		os.Exit(1)
	}

	if err := commit.RunGitAdd(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git add: %v\n", err)
		os.Exit(1)
	}

	config, err := commit.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	if err := commit.RunHook(config.OnBeforeCommit); err != nil {
		fmt.Fprintf(os.Stderr, "Error running onBeforeCommit: %v\n", err)
		os.Exit(1)
	}

	prompt := promptui.Select{
		Label: "Commit type",
		Items: []string{"feat", "fix", "chore", "refactor"},
	}
	_, commitType, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	branchName, err := commit.CurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting branch: %v\n", err)
		os.Exit(1)
	}
	ticket := commit.ExtractJiraTicket(branchName)
	feature := commit.AskFeatureName(ticket)
	message := commit.AskCommitMessage()

	if err := commit.RunHook(config.OnAfterCommit); err != nil {
		fmt.Fprintf(os.Stderr, "Error running onAfterCommit: %v\n", err)
		os.Exit(1)
	}

	fullMessage := commitType + "(" + feature + "): " + message
	fmt.Println("Running: git commit -m \"" + fullMessage + "\"")
	if err := commit.RunGitCommit(fullMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git commit: %v\n", err)
		os.Exit(1)
	}

	if config.PushAfterCommit {
		exists, err := commit.RemoteExists()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check remote: %v\n", err)
			os.Exit(1)
		}
		if !exists {
			commit.PrintRed("\u001f Nessun remote presente. Usa 'git remote add origin <url>' per aggiungerne uno.")
			os.Exit(1)
		}
		fmt.Println("Running: git push")
		if err := commit.RunGitPush(); err != nil {
			fmt.Fprintf(os.Stderr, "Error during git push: %v\n", err)
			os.Exit(1)
		}
	}
}
