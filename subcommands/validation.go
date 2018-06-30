package subcommands

import "github.com/badoux/checkmail"

func validName(name string) bool {
	return len(name) != 0
}

func validJurisdiction(name string) bool {
	// TODO: Implement jurisdiction validation.
	return true
}

func validEMail(email string) bool {
	err := checkmail.ValidateFormat(email)
	return err == nil
}
