package tools

import (
	"fmt"
)

// AskForConfirmation waits for stdin input by the user
// in the cli interface. Input yes or not, then enter (newline).
func AskForConfirmation() bool {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("fatal: ", err)
	}
	positive := []string{"y", "Y", "yes", "Yes", "YES"}
	negative := []string{"n", "N", "no", "No", "NO"}
	if SliceContains(positive, input) {
		return true
	} else if SliceContains(negative, input) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter.")
		return AskForConfirmation()
	}
}
