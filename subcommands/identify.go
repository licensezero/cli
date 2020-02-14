package subcommands

import (
	"flag"
	"io/ioutil"
	"licensezero.com/licensezero/user"
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
		"silent":            silentUsage,
	})

// Identify saves user identification information.
var Identify = &Subcommand{
	Tag:         "misc",
	Description: identifyDescription,
	Handler: func(env Environment) int {
		flagSet := flag.NewFlagSet("identify", flag.ContinueOnError)
		jurisdiction := flagSet.String("jurisdiction", "", "")
		name := flagSet.String("name", "", "")
		email := flagSet.String("email", "", "")
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			env.Stderr.WriteString(identifyUsage)
		}
		err := flagSet.Parse(env.Arguments)
		if err != nil {
			env.Stderr.WriteString("\nError: " + err.Error() + "\n")
			return 1
		}
		if *jurisdiction == "" || *name == "" || *email == "" {
			env.Stderr.WriteString(identifyUsage)
			return 1
		}
		newIdentity := user.Identity{
			Name:         *name,
			Jurisdiction: *jurisdiction,
			EMail:        *email,
		}
		existingIdentity, _ := user.ReadIdentity()
		if existingIdentity != nil && *existingIdentity != newIdentity {
			confirmed, err := env.Stdin.Confirm("Overwrite existing identity?", env.Stdout)
			if err != nil {
				return 1
			}
			if !confirmed {
				return 0
			}
		}
		if !validName(*name) {
			env.Stderr.WriteString("Invalid Name.\n")
			return 1
		}
		if !validJurisdiction(*jurisdiction) {
			env.Stderr.WriteString("Invalid --jurisdiction. Must be ISO 3166-2 code like \"US-CA\" or \"DE-BE\".\n")
			return 1
		}
		if !validEMail(*email) {
			env.Stderr.WriteString("Invalid E-Mail.\n")
			return 1
		}
		err = user.WriteIdentity(&newIdentity)
		if err != nil {
			env.Stderr.WriteString("Could not write identity file.\n")
			return 1
		}
		if !*silent {
			env.Stdout.WriteString("Saved your identification information.\n")
		}
		return 0
	},
}
