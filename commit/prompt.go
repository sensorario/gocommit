package commit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskFeatureName(ticket string) string {
	reader := bufio.NewReader(os.Stdin)
	feature := ""
	if ticket != "" {
		fmt.Printf("The feature name [%s]: ", ticket)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			feature = ticket
		} else {
			feature = input
		}
	} else {
		fmt.Print("The feature name: ")
		input, _ := reader.ReadString('\n')
		feature = strings.TrimSpace(input)
	}
	return feature
}

func AskCommitMessage() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter commit message: ")
	message, _ := reader.ReadString('\n')
	return strings.TrimSpace(message)
}
