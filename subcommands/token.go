package subcommands

import "flag"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const tokenDescription = "Save your API access token."

// Token saves developer IDs and API tokens.
var Token = &Subcommand{
	Tag:         "seller",
	Description: tokenDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("token", flag.ExitOnError)
		developerID := flagSet.String("developer", "", "Developer ID")
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = tokenUsage
		flagSet.Parse(args)
		if *developerID == "" {
			tokenUsage()
		}
		token := secretPrompt("Token: ")
		newDeveloper := data.Developer{
			DeveloperID: *developerID,
			Token:       token,
		}
		existingDeveloper, _ := data.ReadDeveloper(paths.Home)
		if existingDeveloper != nil && *existingDeveloper != newDeveloper {
			if !confirm("Overwrite existing developer info?") {
				os.Exit(0)
			}
		}
		err := data.WriteDeveloper(paths.Home, &newDeveloper)
		if err != nil {
			Fail("Could not write developer file.")
		}
		if !*silent {
			os.Stdout.WriteString("Saved your developer ID and access token.\n")
		}
		os.Exit(0)
	},
}

func tokenUsage() {
	usage := tokenDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero token --developer ID\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"developer ID": "Developer ID (UUID).",
			"silent":       silentLine,
		})
	Fail(usage)
}
