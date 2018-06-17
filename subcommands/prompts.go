package subcommands

import "fmt"
import "golang.org/x/crypto/ssh/terminal"
import "strings"
import "syscall"

func Confirm(prompt string) bool {
	var response string
	fmt.Printf("%s (y/n): ", prompt)
	_, err := fmt.Scan(&response)
	if err != nil {
		panic(err)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	if response == "y" {
		return true
	} else if response == "n" {
		return false
	} else {
		return Confirm(prompt)
	}
}

func SecretPrompt(prompt string) string {
	fmt.Printf(prompt)
	data, err := terminal.ReadPassword(int(syscall.Stderr))
	if err != nil {
		panic(err)
	}
	var response string
	response = string(data)
	fmt.Println()
	return response
}
