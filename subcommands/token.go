package subcommands

import "flag"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const tokenDescription = "Save your API access token."

// Token saves licensor IDs and API tokens.
var Token = &Subcommand{
	Tag:         "seller",
	Description: tokenDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("token", flag.ExitOnError)
		licensorID := flagSet.String("licensor", "", "Licensor ID")
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = tokenUsage
		flagSet.Parse(args)
		if *licensorID == "" {
			tokenUsage()
		}
		token := secretPrompt("Token: ")
		newLicensor := data.Licensor{
			LicensorID: *licensorID,
			Token:      token,
		}
		existingLicensor, _ := data.ReadLicensor(paths.Home)
		if existingLicensor != nil && *existingLicensor != newLicensor {
			if !confirm("Overwrite existing licensor info?") {
				os.Exit(0)
			}
		}
		err := data.WriteLicensor(paths.Home, &newLicensor)
		if err != nil {
			Fail("Could not write licensor file.")
		}
		if !*silent {
			os.Stdout.WriteString("Saved your licensor ID and access token.\n")
		}
		os.Exit(0)
	},
}

func tokenUsage() {
	usage := tokenDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero token --licensor ID\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"licensor ID": "Licensor ID (UUID).",
			"silent":      silentLine,
		})
	Fail(usage)
}
