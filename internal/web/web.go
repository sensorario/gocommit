package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gocommit/commit"
)

const addr = "http://localhost:8080"
const pidFile = "/tmp/gocommit-web.pid"

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
	pid := cmd.Process.Pid
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not write pid file: %v\n", err)
	}
	commit.PrintGreen("Web server started at " + addr + " (pid " + fmt.Sprint(pid) + ")")
	exec.Command("open", addr).Start()
}

func Kill() {
	pid, err := pidByPort("8080")
	if err != nil || pid == 0 {
		// fall back to pid file
		data, ferr := os.ReadFile(pidFile)
		if ferr != nil {
			commit.PrintRed("Web server does not appear to be running on port 8080.")
			os.Exit(1)
		}
		pid, err = strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			commit.PrintRed("Invalid pid file content.")
			os.Exit(1)
		}
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		commit.PrintRed(fmt.Sprintf("Could not find process %d: %v", pid, err))
		os.Exit(1)
	}
	if err := proc.Kill(); err != nil {
		commit.PrintRed(fmt.Sprintf("Could not kill process %d: %v", pid, err))
		os.Exit(1)
	}
	os.Remove(pidFile)
	commit.PrintGreen(fmt.Sprintf("Web server (pid %d) stopped.", pid))
}

// pidByPort returns the PID of the process listening on the given TCP port.
// Uses lsof, which is available on macOS and most Linux systems.
func pidByPort(port string) (int, error) {
	out, err := exec.Command("lsof", "-ti", "tcp:"+port).Output()
	if err != nil {
		return 0, err
	}
	pidStr := strings.TrimSpace(string(out))
	if pidStr == "" {
		return 0, nil
	}
	// lsof may return multiple lines; take the first
	lines := strings.Split(pidStr, "\n")
	return strconv.Atoi(strings.TrimSpace(lines[0]))
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
			   status := http.StatusInternalServerError
			   if _, ok := err.(*commit.CheckoutError); ok {
				   status = http.StatusConflict
			   }
			   http.Error(w, "checkout failed: "+err.Error(), status)
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
