package subcommands

import "flag"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

const identifyDescription = "Save your identity information."

// Identify saves user identification information.
var Identify = &Subcommand{
	Tag:         "misc",
	Description: identifyDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("identify", flag.ExitOnError)
		jurisdiction := flagSet.String("jurisdiction", "", "")
		name := flagSet.String("name", "", "")
		email := flagSet.String("email", "", "")
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = identifyUsage
		flagSet.Parse(args)
		if *jurisdiction == "" || *name == "" || *email == "" {
			identifyUsage()
		}
		newIdentity := data.Identity{
			Name:         *name,
			Jurisdiction: *jurisdiction,
			EMail:        *email,
		}
		existingIdentity, _ := data.ReadIdentity(paths.Home)
		if existingIdentity != nil && *existingIdentity != newIdentity {
			if !confirm("Overwrite existing identity?") {
				os.Exit(0)
			}
		}
		if !validName(*name) {
			Fail("Invalid Name.")
		}
		if !validJurisdiction(*jurisdiction) {
			Fail("Invalid Jurisdiction.")
		}
		if !validEMail(*email) {
			Fail("Invalid E-Mail.")
		}
		err := data.WriteIdentity(paths.Home, &newIdentity)
		if err != nil {
			Fail("Could not write identity file.")
		}
		if !*silent {
			os.Stdout.WriteString("Saved your identification information.\n")
		}
		os.Exit(0)
	},
}

func identifyUsage() {
	usage := identifyDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero identify --name NAME --jurisdiction CODE --email ADDRESS\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"email ADDRESS":     "Your e-mail address",
			"jurisdiction CODE": "Your tax jurisdiction (ISO 3166-2, like \"US-CA\")",
			"name NAME":         "Your full name.",
			"silent":            silentLine,
		})
	Fail(usage)
}
