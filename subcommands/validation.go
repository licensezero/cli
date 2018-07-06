package subcommands

import "regexp"
import "github.com/badoux/checkmail"

func validName(name string) bool {
	return len(name) != 0
}

var jurisdictionPattern = regexp.MustCompile(`^[A-Z]{2}-[0-9A-Z]{1,3}$`)

func validJurisdiction(name string) bool {
	return jurisdictionPattern.MatchString(name)
}

func validEMail(email string) bool {
	err := checkmail.ValidateFormat(email)
	return err == nil
}
