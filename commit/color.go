package commit

import "fmt"

func PrintRed(msg string) {
	fmt.Printf("\033[31m%s\033[0m\n", msg)
}
