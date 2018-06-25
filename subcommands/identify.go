package subcommands

import "flag"
import "github.com/licensezero/cli/data"
import "os"

const identifyDescription = "Wave your personal details for quoting and buying licenses."

var Identify = Subcommand{
	Description: identifyDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("identify", flag.ExitOnError)
		jurisdiction := flagSet.String("jurisdiction", "", "")
		name := flagSet.String("name", "", "")
		email := flagSet.String("email", "", "")
		silent := Silent(flagSet)
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
			if !Confirm("Overwrite existing identity?") {
				os.Exit(0)
			}
		}
		if !ValidName(*name) {
			os.Stderr.WriteString("Invalid Name.\n")
			os.Exit(1)
		}
		if !ValidJurisdiction(*jurisdiction) {
			os.Stderr.WriteString("Invalid Jurisdiction.\n")
			os.Exit(1)
		}
		if !ValidEMail(*email) {
			os.Stderr.WriteString("Invalid E-Mail.\n")
			os.Exit(1)
		}
		err := data.WriteIdentity(paths.Home, &newIdentity)
		if err != nil {
			os.Stderr.WriteString("Could not write identity file.\n")
			os.Exit(1)
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
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
