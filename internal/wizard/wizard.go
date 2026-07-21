package wizard

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"

	"gocommit/commit"
)

const customFeatureOption = "Custom..."

var CommitTypes = []string{
	"chore",
	"fix",
	"test",
	"feat",
	"refactor",
}

func selectFeatureName(ticket string) string {
	recent := commit.RecentFeatureNames(5)
	items := make([]string, 0, len(recent)+1)
	items = append(items, recent...)
	items = append(items, customFeatureOption)

	sel := promptui.Select{
		Label: "Feature name",
		Items: items,
	}
	_, chosen, err := sel.Run()
	if err != nil || chosen == customFeatureOption {
		return commit.AskFeatureName(ticket)
	}
	return chosen
}

func Run() {
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
		Items: CommitTypes,
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

	feature := selectFeatureName(ticket)
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
			commit.PrintRed("Nessun remote presente. Usa 'git remote add origin <url>' per aggiungerne uno.")
			os.Exit(1)
		}

		hasUpstream, err := commit.BranchHasUpstream()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check upstream: %v\n", err)
			os.Exit(1)
		}

		if !hasUpstream {
			sel := promptui.Select{
				Label: "Il branch non ha un upstream remoto. Impostare 'origin/" + branchName + "' come upstream?",
				Items: []string{"Si", "No"},
			}
			_, choice, err := sel.Run()
			if err != nil || choice == "No" {
				fmt.Println("Push saltato.")
				return
			}
			fmt.Println("Running: git push --set-upstream origin " + branchName)
			if err := commit.RunGitPushSetUpstream(branchName); err != nil {
				fmt.Fprintf(os.Stderr, "Error during git push --set-upstream: %v\n", err)
				os.Exit(1)
			}
			return
		}

		fmt.Println("Running: git push")
		if err := commit.RunGitPush(); err != nil {
			fmt.Fprintf(os.Stderr, "Error during git push: %v\n", err)
			os.Exit(1)
		}
	}
}