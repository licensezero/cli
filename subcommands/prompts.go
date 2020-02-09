package subcommands

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"licensezero.com/licensezero/api"
	"os"
	"strings"
)

func confirm(prompt string, stdin, stdout *os.File) (bool, error) {
	var response string
	stdout.WriteString(prompt + " (y/n): ")
	_, err := fmt.Scan(stdin, &response)
	if err != nil {
		return false, err
	}
	response = strings.TrimSpace(strings.ToLower(response))
	if response == "y" {
		return true, nil
	} else if response == "n" {
		return false, nil
	} else {
		return confirm(prompt, stdin, stdout)
	}
}

func secretPrompt(prompt string, stdin, stdout *os.File) (response string, err error) {
	stdout.WriteString(prompt)
	data, err := terminal.ReadPassword(int(stdin.Fd()))
	if err != nil {
		return
	}
	response = string(data)
	stdout.WriteString("\n")
	return
}

const termsPrompt = "Do you agree to " + api.TermsReference + "?"

func confirmTermsOfService(stdin, stdout *os.File) (bool, error) {
	return confirm(termsPrompt, stdin, stdout)
}

const agencyPrompt = "Do you agree to " + api.AgencyReference + "?"

func confirmAgencyTerms(stdin, stdout *os.File) (bool, error) {
	return confirm(agencyPrompt, stdin, stdout)
}
