package subcommands

import (
	"github.com/badoux/checkmail"
	"licensezero.com/licensezero/api"
	"regexp"
)

func validJurisdiction(j string) bool {
	return api.ValidateJurisdiction(j)
}

func validName(name string) bool {
	return len(name) != 0
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

func invalidJurisdiction() {
	Fail("Invalid --jurisdiction. Must be ISO 3166-2 code like \"US-CA\" or \"DE-BE\".")
}
