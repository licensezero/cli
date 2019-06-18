package subcommands

import "regexp"
import "github.com/badoux/checkmail"

func validName(name string) bool {
	return len(name) != 0
}

//go:generate ../generate-jurisdictions

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
