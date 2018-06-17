package subcommands

import "os"

var Identify = Subcommand{
	Description: "Identify yourself.",
	Handler: func(args []string, paths Paths) {
		if len(args) != 3 {
			os.Stderr.WriteString("<name> <jurisdiction> <email>\n")
			os.Exit(1)
		} else {
			name := args[0]
			jurisdiction := args[1]
			email := args[2]
			newIdentity := Identity{
				Name:         name,
				Jurisdiction: jurisdiction,
				EMail:        email,
			}
			existingIdentity, _ := readIdentity(paths.Home)
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
			err := writeIdentity(paths.Home, &newIdentity)
			if err != nil {
				os.Stderr.WriteString("Could not write identity file.\n")
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}
	},
}
