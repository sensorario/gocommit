package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gocommit/commit"
)

var shutdownCh = make(chan struct{}, 1)

const defaultPort = 8080

type serverInfo struct {
	PID  int `json:"pid"`
	Port int  `json:"port"`
}

func qweDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp"
	}
	dir := home + "/.sensorario-qwe"
	os.MkdirAll(dir, 0700)
	return dir
}

func serverInfoPath() string {
	return qweDir() + "/server.json"
}

func serverAddr(port int) string {
	return fmt.Sprintf("http://localhost:%d", port)
}

func loadServerInfo() (serverInfo, bool) {
	data, err := os.ReadFile(serverInfoPath())
	if err != nil {
		return serverInfo{}, false
	}
	var info serverInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return serverInfo{}, false
	}
	return info, info.PID > 0 && info.Port > 0
}

func saveServerInfo(info serverInfo) {
	data, _ := json.Marshal(info)
	if err := os.WriteFile(serverInfoPath(), data, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not write server.json: %v\n", err)
	}
}

func isPortInUse(port int) bool {
	pid, err := pidByPort(strconv.Itoa(port))
	return err == nil && pid != 0
}

func findFreePort(start int) int {
	for port := start; port < start+100; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port
		}
	}
	return start
}

func Run() {
	if info, ok := loadServerInfo(); ok && isPortInUse(info.Port) {
		addr := serverAddr(info.Port)
		commit.PrintGreen("Web server already running at " + addr + " (pid " + fmt.Sprint(info.PID) + ")")
		exec.Command("open", addr).Start()
		return
	}

	port := findFreePort(defaultPort)
	addr := serverAddr(port)

	self, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable: %v\n", err)
		os.Exit(1)
	}
	cmd := exec.Command(self, "--web-server")
	cmd.Env = append(os.Environ(), fmt.Sprintf("QWE_PORT=%d", port))
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting web server: %v\n", err)
		os.Exit(1)
	}

	saveServerInfo(serverInfo{PID: cmd.Process.Pid, Port: port})
	commit.PrintGreen("Web server started at " + addr + " (pid " + fmt.Sprint(cmd.Process.Pid) + ")")
	exec.Command("open", addr).Start()
}

func Kill() {
	info, ok := loadServerInfo()
	if !ok {
		commit.PrintRed("Web server does not appear to be running (no server.json found).")
		os.Exit(1)
	}

	proc, err := os.FindProcess(info.PID)
	if err != nil {
		commit.PrintRed(fmt.Sprintf("Could not find process %d: %v", info.PID, err))
		os.Remove(serverInfoPath())
		os.Exit(1)
	}
	if err := proc.Kill(); err != nil {
		commit.PrintRed(fmt.Sprintf("Could not kill process %d: %v", info.PID, err))
		os.Exit(1)
	}
	os.Remove(serverInfoPath())
	commit.PrintGreen(fmt.Sprintf("Web server (pid %d, port %d) stopped.", info.PID, info.Port))
}

func pidByPort(port string) (int, error) {
	out, err := exec.Command("lsof", "-ti", "tcp:"+port).Output()
	if err != nil {
		return 0, err
	}
	pidStr := strings.TrimSpace(string(out))
	if pidStr == "" {
		return 0, nil
	}
	lines := strings.Split(pidStr, "\n")
	return strconv.Atoi(strings.TrimSpace(lines[0]))
}

func Serve() {
	RegisterCurrentRepo()
	RegisterRoutes()

	port := defaultPort
	if p, err := strconv.Atoi(os.Getenv("QWE_PORT")); err == nil && p > 0 {
		port = p
	}

	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	go func() {
		<-shutdownCh
		srv.Shutdown(context.Background())
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		os.Exit(1)
	}
}