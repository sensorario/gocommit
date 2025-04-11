package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	fmt.Println("Running: git add .")
	if err := runCommand("git", "add", "."); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git add: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter commit message: ")
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading message: %v\n", err)
		os.Exit(1)
	}
	message = strings.TrimSpace(message)

	fmt.Println("Running: git commit -m \"" + message + "\"")
	if err := runCommand("git", "commit", "-m", message); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git commit: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Running: git push")
	if err := runCommand("git", "push"); err != nil {
		fmt.Fprintf(os.Stderr, "Error during git push: %v\n", err)
		os.Exit(1)
	}
}