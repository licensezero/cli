package subcommands

import (
	"flag"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/user"
	"net/http"
)

const identifyDescription = "Save your identity information."

var identifyUsage = identifyDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero identify --name NAME --jurisdiction CODE --email ADDRESS\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"email ADDRESS":     "Your e-mail address",
		"jurisdiction CODE": "Your tax jurisdiction (ISO 3166-2, like \"US-CA\")",
		"name NAME":         "Your full name.",
		"silent":            silentLine,
	})

// Identify saves user identification information.
var Identify = &Subcommand{
	Tag:         "misc",
	Description: identifyDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		flagSet := flag.NewFlagSet("identify", flag.ContinueOnError)
		jurisdiction := flagSet.String("jurisdiction", "", "")
		name := flagSet.String("name", "", "")
		email := flagSet.String("email", "", "")
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			stderr.WriteString(identifyUsage)
		}
		err := flagSet.Parse(args)
		if err != nil {
			stderr.WriteString("\nError: " + err.Error() + "\n")
			return 1
		}
		if *jurisdiction == "" || *name == "" || *email == "" {
			stderr.WriteString(identifyUsage)
			return 1
		}
		newIdentity := user.Identity{
			Name:         *name,
			Jurisdiction: *jurisdiction,
			EMail:        *email,
		}
		existingIdentity, _ := user.ReadIdentity()
		if existingIdentity != nil && *existingIdentity != newIdentity {
			confirmed, err := stdin.Confirm("Overwrite existing identity?", stdout)
			if err != nil {
				return 1
			}
			if !confirmed {
				return 0
			}
		}
		if !validName(*name) {
			stderr.WriteString("Invalid Name.\n")
			return 1
		}
		if !validJurisdiction(*jurisdiction) {
			stderr.WriteString("Invalid --jurisdiction. Must be ISO 3166-2 code like \"US-CA\" or \"DE-BE\".\n")
			return 1
		}
		if !validEMail(*email) {
			stderr.WriteString("Invalid E-Mail.\n")
			return 1
		}
		err = user.WriteIdentity(&newIdentity)
		if err != nil {
			stderr.WriteString("Could not write identity file.\n")
			return 1
		}
		if !*silent {
			stdout.WriteString("Saved your identification information.\n")
		}
		return 0
	},
}
