package web

import (
	"os"
)

func GetVersion() string {
	data, err := os.ReadFile("VERSION")
	if err != nil {
		return "unknown"
	}
	return string(data)
}
