package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

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
			fmt.Printf("  %-20s %s\n", "branch", "Interactively switch to another git branch")
			fmt.Printf("  %-20s %s\n", "check", "Verify gocommit.conf.json exists and all variables are present; adds missing ones with defaults")
			fmt.Printf("  %-20s %s\n", "web", "Start a local web server and open it in the browser")
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
		if arg == "branch" {
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
			return
		}
		if arg == "web" {
			addr := "http://localhost:8080"
			http.HandleFunc("/branches", func(w http.ResponseWriter, r *http.Request) {
				branches, err := commit.ListLocalBranches()
				if err != nil {
					http.Error(w, "failed to list branches", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(branches)
			})
			http.HandleFunc("/checkout", func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}
				var body struct {
					Branch string `json:"branch"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Branch == "" {
					http.Error(w, "invalid request", http.StatusBadRequest)
					return
				}
				if err := commit.CheckoutBranch(body.Branch); err != nil {
					http.Error(w, "checkout failed: "+err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"branch": body.Branch})
			})
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>gocommit</title>
  <style>
    body { font-family: sans-serif; padding: 2rem; }
    button { display: block; margin: 0.25rem 0; padding: 0.4rem 1rem; cursor: pointer; }
    #message { margin-top: 1rem; font-weight: bold; }
  </style>
</head>
<body>
  <h1>Hello World</h1>
  <h2>Branches</h2>
  <div id="branches"></div>
  <div id="message"></div>
  <script>
    function checkout(branch) {
      fetch('/checkout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ branch })
      })
      .then(r => r.json())
      .then(data => {
        document.getElementById('message').textContent = 'Switched to: ' + data.branch;
      })
      .catch(() => {
        document.getElementById('message').textContent = 'Error switching branch.';
      });
    }

    fetch('/branches')
      .then(r => r.json())
      .then(branches => {
        const container = document.getElementById('branches');
        branches.forEach(b => {
          const btn = document.createElement('button');
          btn.textContent = b;
          btn.onclick = () => checkout(b);
          container.appendChild(btn);
        });
      });
  </script>
</body>
</html>`)
			})
			fmt.Println("Serving on " + addr)
			exec.Command("open", addr).Start()
			if err := http.ListenAndServe(":8080", nil); err != nil {
				fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
				os.Exit(1)
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
