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

func validID(uuid string) bool {
	// UUIDv4
	hex := "[a-f0-9]"
	re := regexp.MustCompile("^" + hex + "{8}" + "-" + hex + "{4}" + "-" + "4" + hex + "{3}" + "-" + "[8|9|a|b]" + hex + "{3}-" + hex + "{12}" + "$")
	return re.MatchString(uuid)
}

func invalidID() {
	Fail("Invalid --id. Must be UUID from `licensezero offer`.")
}
