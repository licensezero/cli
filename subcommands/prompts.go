package subcommands

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"licensezero.com/licensezero/api"
	"os"
	"strings"
)

// InputDevice abstacts over stdin so we can mock it in tests.
type InputDevice interface {
	Confirm(prompt string, stdout io.StringWriter) (bool, error)
	SecretPrompt(prompt string, stdout io.StringWriter) (string, error)
}

// StandardInputDevice is an InputDevice based on an actual file.
type StandardInputDevice struct {
	File *os.File
}

// Confirm prompts by scanning d.File.
func (d *StandardInputDevice) Confirm(prompt string, stdout io.StringWriter) (bool, error) {
	var response string
	stdout.WriteString(prompt + " (y/n): ")
	_, err := fmt.Scan(d.File, &response)
	if err != nil {
		return false, err
	}
	response = strings.TrimSpace(strings.ToLower(response))
	if response == "y" {
		return true, nil
	} else if response == "n" {
		return false, nil
	} else {
		return d.Confirm(prompt, stdout)
	}
}

// SecretPrompt prompts with terminal.ReadPassword.
func (d *StandardInputDevice) SecretPrompt(prompt string, stdout io.StringWriter) (response string, err error) {
	stdout.WriteString(prompt)
	data, err := terminal.ReadPassword(int(d.File.Fd()))
	if err != nil {
		return
	}
	response = string(data)
	stdout.WriteString("\n")
	return
}

const termsPrompt = "Do you agree to " + api.TermsReference + "?"

func confirmTermsOfService(input InputDevice, stdout io.StringWriter) (bool, error) {
	return input.Confirm(termsPrompt, stdout)
}

const agencyPrompt = "Do you agree to " + api.AgencyReference + "?"

func confirmAgencyTerms(input InputDevice, stdout io.StringWriter) (bool, error) {
	return input.Confirm(agencyPrompt, stdout)
}
