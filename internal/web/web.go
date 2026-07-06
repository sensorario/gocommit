package web

import (
	"context"
	"fmt"
	"gocommit/commit"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var shutdownCh = make(chan struct{}, 1)

const addr = "http://localhost:8080"

func pidFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/gocommit-web.pid"
	}
	dir := home + "/.sensorario-qwe"
	os.MkdirAll(dir, 0700)
	return dir + "/gocommit-web.pid"
}

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
	if err := os.WriteFile(pidFilePath(), []byte(strconv.Itoa(pid)), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not write pid file: %v\n", err)
	}
	commit.PrintGreen("Web server started at " + addr + " (pid " + fmt.Sprint(pid) + ")")
	exec.Command("open", addr).Start()
}

func Kill() {
	pid, err := pidByPort("8080")
	if err != nil || pid == 0 {
		// fall back to pid file
		data, ferr := os.ReadFile(pidFilePath())
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
	os.Remove(pidFilePath())
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
	RegisterRoutes()
	srv := &http.Server{Addr: ":8080"}
	go func() {
		<-shutdownCh
		srv.Shutdown(context.Background())
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		os.Exit(1)
	}
}
