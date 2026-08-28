package main

import (
	"fmt"
	"os"

	"gocommit/commit"
	"gocommit/internal/branch"
	"gocommit/internal/check"
	"gocommit/internal/versions"
	"gocommit/internal/web"
	"gocommit/internal/wizard"
)


var Version = EmbeddedVersion

func init() {
	web.Version = Version
}

func main() {
	// internal flag used by 'qwe web' to start the background HTTP server
	if len(os.Args) > 1 && os.Args[1] == "--web-server" {
		web.Serve()
		return
	}

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
			fmt.Printf("  %-20s %s\n", "branch", "Interactively switch to another git branch")
			fmt.Printf("  %-20s %s\n", "check", "Verify gocommit.conf.json exists and all variables are present; adds missing ones with defaults")
			fmt.Printf("  %-20s %s\n", "web", "Start a local web server in the background and open it in the browser")
			fmt.Printf("  %-20s %s\n", "web-stop", "Stop the background web server")
			fmt.Printf("  %-20s %s\n", "versions", "Show the package.json version of every known repository")
			fmt.Printf("  %-20s %s\n", "help, -h, --help", "Show this help message")
			fmt.Printf("  %-20s %s\n", "-v, --version", "Print the current version")
			return
		}
		commands := map[string]func(){
			"branch":   branch.Run,
			"check":    check.Run,
			"web":      web.Run,
			"web-stop": web.Kill,
			"versions": versions.Run,
		}
		if fn, ok := commands[arg]; ok {
			fn()
			return
		}
		commit.PrintRed("Unknown command: " + arg + ". Run 'qwe help' to see available commands.")
		os.Exit(1)
	}
	wizard.Run()
}

