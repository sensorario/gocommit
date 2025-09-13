package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
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
		printRed(" Non sono accettati parametri ad eccezione di -v o --version")
		os.Exit(1)
	}

	// Ensure gocommit.conf.json exists in the project root
	configPath := "gocommit.conf.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		f, err := os.Create(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", configPath, err)
			os.Exit(1)
		}
		config := map[string]string{
			"onBeforeCommit": "",
			"onAfterCommit":  "",
		}
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to %s: %v\n", configPath, err)
			f.Close()
			os.Exit(1)
		}
		f.Close()
	}

	fmt.Println("Running: git add .")
	if err := runCommand("git", "add", "."); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git add: %v\n", err)
		os.Exit(1)
	}

	// Check for gocommit.conf.json and run onBeforeCommit if set
	confPath := "gocommit.conf.json"
	if _, err := os.Stat(confPath); err == nil {
		confFile, err := os.Open(confPath)
		if err == nil {
			defer confFile.Close()
			var onBeforeCommit string
			buf := make([]byte, 4096)
			n, _ := confFile.Read(buf)
			confStr := string(buf[:n])

			// Simple parse for onBeforeCommit value
			idx := strings.Index(confStr, "\"onBeforeCommit\"")
			if idx != -1 {
				rest := confStr[idx+len("\"onBeforeCommit\""):]
				colonIdx := strings.Index(rest, ":")
				if colonIdx != -1 {
					rest = rest[colonIdx+1:]
					quoteIdx := strings.Index(rest, "\"")
					if quoteIdx != -1 {
						rest = rest[quoteIdx+1:]
						endQuoteIdx := strings.Index(rest, "\"")
						if endQuoteIdx != -1 {
							onBeforeCommit = rest[:endQuoteIdx]
						}
					}
				}
			}

			if len(onBeforeCommit) > 0 {
				fmt.Printf("Running onBeforeCommit: %s\n", onBeforeCommit)
				// Run the command using shell
				shellCmd := exec.Command("sh", "-c", onBeforeCommit)
				shellCmd.Stdout = os.Stdout
				shellCmd.Stderr = os.Stderr
				if err := shellCmd.Run(); err != nil {
					fmt.Fprintf(os.Stderr, "Error running onBeforeCommit: %v\n", err)
					os.Exit(1)
				}
			}
		}
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

	// Check for onAfterCommit in configuration and run it after commit
	if _, err := os.Stat(confPath); err == nil {
		confFile, err := os.Open(confPath)
		if err == nil {
			defer confFile.Close()
			var onAfterCommit string
			buf := make([]byte, 4096)
			n, _ := confFile.Read(buf)
			confStr := string(buf[:n])

			// Simple parse for onAfterCommit value
			idx := strings.Index(confStr, "\"onAfterCommit\"")
			if idx != -1 {
				rest := confStr[idx+len("\"onAfterCommit\""):]
				colonIdx := strings.Index(rest, ":")
				if colonIdx != -1 {
					rest = rest[colonIdx+1:]
					quoteIdx := strings.Index(rest, "\"")
					if quoteIdx != -1 {
						rest = rest[quoteIdx+1:]
						endQuoteIdx := strings.Index(rest, "\"")
						if endQuoteIdx != -1 {
							onAfterCommit = rest[:endQuoteIdx]
						}
					}
				}
			}

			if len(onAfterCommit) > 0 {
				fmt.Printf("Running onAfterCommit: %s\n", onAfterCommit)
				shellCmd := exec.Command("sh", "-c", onAfterCommit)
				shellCmd.Stdout = os.Stdout
				shellCmd.Stderr = os.Stderr
				if err := shellCmd.Run(); err != nil {
					fmt.Fprintf(os.Stderr, "Error running onAfterCommit: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

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
		printRed(" Nessun remote presente. Usa 'git remote add origin <url>' per aggiungerne uno.")
		os.Exit(1)
	}

	fmt.Println("Running: git push")
	if err := runCommand("git", "push"); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git push: %v\n", err)
		os.Exit(1)
	}
}
