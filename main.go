package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitRemoteExists() (bool, error) {
	cmd := exec.Command("git", "remote")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(output)) != "", nil
}

func printRed(msg string) {
	fmt.Printf("\033[31m%s\033[0m\n", msg) // ANSI escape code for red
}

func main() {
	fmt.Println("Running: git add .")
	if err := runCommand("git", "add", "."); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git add: %v\n", err)
		os.Exit(1)
	}

	prompt := promptui.Select{
		Label: "Commit type",
		Items: []string{"wip", "feat", "fix", "chore", "docs", "style", "refactor", "perf", "test"},
	}
	_, commitType, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("The feature name: ")
	feature, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading message: %v\n", err)
		os.Exit(1)
	}
	feature = strings.TrimSpace(feature)

	fmt.Print("Enter commit message: ")
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading message: %v\n", err)
		os.Exit(1)
	}
	message = strings.TrimSpace(message)

	fullMessage := commitType + "(" + feature + "): " + message
	fmt.Println("Running: git commit -m \"" + fullMessage + "\"")
	if err := runCommand("git", "commit", "-m", fullMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git commit: %v\n", err)
		os.Exit(1)
	}

	exists, err := gitRemoteExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check remote: %v\n", err)
		os.Exit(1)
	}

	if !exists {
		printRed("❌ Nessun remote presente. Usa 'git remote add origin <url>' per aggiungerne uno.")
		os.Exit(1)
	}

	fmt.Println("Running: git push")
	if err := runCommand("git", "push"); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git push: %v\n", err)
		os.Exit(1)
	}
}
