package web

import "os"

// Version viene impostata da main.go
var Version string

func GetVersion() string {
	if os.Getenv("DEV_MODE") == "1" {
		data, err := os.ReadFile("VERSION")
		if err == nil {
			return string(data)
		}
	}
	return Version
}
