package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"gocommit/commit"
)

const addr = "http://localhost:8080"

func Run() {
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable: %v\n", err)
		os.Exit(1)
	}
	cmd := exec.Command(self, "--web-server")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting web server: %v\n", err)
		os.Exit(1)
	}
	commit.PrintGreen("Web server started at " + addr + " (pid " + fmt.Sprint(cmd.Process.Pid) + ")")
	exec.Command("open", addr).Start()
}

func Serve() {
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
      .then(r => {
        if (!r.ok) {
          return r.text().then(t => { throw new Error(t); });
        }
        return r.json();
      })
      .then(data => {
        document.getElementById('message').textContent = 'Switched to: ' + data.branch;
      })
      .catch(err => {
        document.getElementById('message').textContent = 'Error: ' + err.message;
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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		os.Exit(1)
	}
}
