package subcommands

import "flag"
import "github.com/licensezero/cli/data"
import "os"

const identifyDescription = "Identify yourself for buying and checking licenses."

var Identify = Subcommand{
	Description: identifyDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("identify", flag.ExitOnError)
		jurisdiction := flagSet.String("jurisdiction", "", "")
		name := flagSet.String("name", "", "")
		email := flagSet.String("email", "", "")
		flagSet.Usage = identifyUsage
		flagSet.Parse(args)
		if len(*jurisdiction) == 0 || len(*name) == 0 || len(*email) == 0 {
			identifyUsage()
		}
		if len(args) != 3 {
			identifyUsage()
		} else {
			name := args[0]
			jurisdiction := args[1]
			email := args[2]
			newIdentity := data.Identity{
				Name:         name,
				Jurisdiction: jurisdiction,
				EMail:        email,
			}
			existingIdentity, _ := data.ReadIdentity(paths.Home)
			if existingIdentity != nil && *existingIdentity != newIdentity {
				if !Confirm("Overwrite existing identity?") {
					os.Exit(0)
				}
			}
			if !ValidName(name) {
				os.Stderr.WriteString("Invalid Name.\n")
				os.Exit(1)
			}
			if !ValidJurisdiction(jurisdiction) {
				os.Stderr.WriteString("Invalid Jurisdiction.\n")
				os.Exit(1)
			}
			if !ValidEMail(email) {
				os.Stderr.WriteString("Invalid E-Mail.\n")
				os.Exit(1)
			}
			err := data.WriteIdentity(paths.Home, &newIdentity)
			if err != nil {
				os.Stderr.WriteString("Could not write identity file.\n")
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}
	},
}

func identifyUsage() {
	usage := identifyDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero identify --name NAME --jurisdiction CODE --email ADDRESS\n\n" +
		"Options:\n" +
		"  --email ADDRESS      Your e-mail address\n" +
		"  --jurisdiction CODE  Your tax jurisdiction (ISO 3166-2, like \"US-CA\")\n" +
		"  --name NAME          Your full name.\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
