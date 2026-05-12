package commit

import "fmt"

func PrintRed(msg string) {
	fmt.Printf("\033[31m%s\033[0m\n", msg)
}

func PrintGreen(msg string) {
	fmt.Printf("\033[32m%s\033[0m\n", msg)
}

func PrintYellow(msg string) {
	fmt.Printf("\033[33m%s\033[0m\n", msg)
}
