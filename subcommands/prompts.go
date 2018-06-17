package subcommands

import "fmt"
import "strings"

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
