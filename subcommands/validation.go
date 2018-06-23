package subcommands

import "github.com/badoux/checkmail"

func ValidName(name string) bool {
	return len(name) != 0
}

func ValidJurisdiction(name string) bool {
	// TODO: Implement jurisdiction validation.
	return true
}

func ValidEMail(email string) bool {
	err := checkmail.ValidateFormat(email)
	return err == nil
}
