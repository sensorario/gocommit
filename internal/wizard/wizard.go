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

// resolveDefaultBranch returns "next" if it exists locally; otherwise it asks
// the user to pick the main branch from the list of local branches.
func resolveDefaultBranch() string {
	branches, err := commit.ListLocalBranches()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing branches: %v\n", err)
		os.Exit(1)
	}
	for _, b := range branches {
		if b == "next" {
			return "next"
		}
	}

	sel := promptui.Select{
		Label: "Nessun branch 'next' trovato. Seleziona il branch principale",
		Items: branches,
	}
	_, chosen, err := sel.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed: %v\n", err)
		os.Exit(1)
	}
	return chosen
}

func mergeIntoNext(branchName, defaultBranch string) {
	fmt.Println("Running: git checkout " + defaultBranch)
	if err := commit.CheckoutBranch(defaultBranch); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git checkout %s: %v\n", defaultBranch, err)
		os.Exit(1)
	}

	fmt.Println("Running: git merge --no-ff " + branchName)
	if err := commit.RunGitMerge(branchName); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git merge: %v\n", err)
		os.Exit(1)
	}
	commit.PrintGreen("Branch '" + branchName + "' mergiato su " + defaultBranch + ".")

	sel := promptui.Select{
		Label: "Cancellare il branch '" + branchName + "' ora mergiato (locale ed eventuale remoto)?",
		Items: []string{"Si", "No"},
	}
	_, choice, err := sel.Run()
	if err != nil || choice == "No" {
		return
	}

	if err := commit.DeleteLocalBranch(branchName); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git branch -d: %v\n", err)
		os.Exit(1)
	}
	commit.PrintGreen("Branch locale '" + branchName + "' cancellato.")

	remoteExists, err := commit.RemoteBranchExists(branchName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking remote branch: %v\n", err)
		os.Exit(1)
	}
	if remoteExists {
		if err := commit.DeleteRemoteBranch(branchName); err != nil {
			fmt.Fprintf(os.Stderr, "Error during git push origin --delete: %v\n", err)
			os.Exit(1)
		}
		commit.PrintGreen("Branch remoto '" + branchName + "' cancellato.")
	}
}

func Run() {
	configPath := "gocommit.conf.json"
	if err := commit.EnsureConfig(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config: %v\n", err)
		os.Exit(1)
	}

	defaultBranch := resolveDefaultBranch()

	files, err := commit.GetUncommittedFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting git status: %v\n", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Println("Nothing to commit.")

		branchName, err := commit.CurrentBranch()
		if err == nil && branchName != defaultBranch {
			sel := promptui.Select{
				Label: "Nessuna modifica. Mergiare '" + branchName + "' su " + defaultBranch + " e cancellare il branch?",
				Items: []string{"Si", "No"},
			}
			_, choice, err := sel.Run()
			if err == nil && choice == "Si" {
				mergeIntoNext(branchName, defaultBranch)
			}
		}
		os.Exit(0)
	}
	fmt.Println("Files to be committed:")
	for _, f := range files {
		fmt.Println(" ", f)
	}
	fmt.Println()

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
			mergeChoice := "Merge su " + defaultBranch
			sel := promptui.Select{
				Label: "Il branch non ha un upstream remoto. Cosa vuoi fare?",
				Items: []string{"Push (imposta upstream)", mergeChoice},
			}
			_, choice, err := sel.Run()
			if err != nil {
				fmt.Println("Push saltato.")
				return
			}

			if choice == mergeChoice {
				mergeIntoNext(branchName, defaultBranch)
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