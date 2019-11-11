package subcommands

import "fmt"
import "golang.org/x/crypto/ssh/terminal"
import "licensezero.com/cli/api"
import "os"
import "strings"

func confirm(prompt string) bool {
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
		return confirm(prompt)
	}
}

func secretPrompt(prompt string) string {
	fmt.Printf(prompt)
	data, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	response := string(data)
	fmt.Println()
	return response
}

const termsPrompt = "Do you agree to " + api.TermsReference + "?"

func confirmTermsOfService() bool {
	return confirm(termsPrompt)
}

const agencyPrompt = "Do you agree to " + api.AgencyReference + "?"

func confirmAgencyTerms() bool {
	return confirm(agencyPrompt)
}
