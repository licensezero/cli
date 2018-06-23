package subcommands

import "fmt"
import "github.com/licensezero/cli/data"
import "os"

var SetLicensorID = Subcommand{
	Description: "Set your licensezero.com licensor ID.",
	Handler: func(args []string, paths Paths) {
		if len(args) != 1 {
			os.Stderr.WriteString("Usage:\n\t<licensor ID>\n")
			os.Exit(1)
		} else {
			licensorID := args[0]
			fmt.Println(licensorID)
			token := SecretPrompt("Token: ")
			newLicensor := data.Licensor{
				LicensorID: licensorID,
				Token:      token,
			}
			existingLicensor, _ := data.ReadLicensor(paths.Home)
			if existingLicensor != nil && *existingLicensor != newLicensor {
				if !Confirm("Overwrite existing licensor info?") {
					os.Exit(0)
				}
			}
			err := data.WriteLicensor(paths.Home, &newLicensor)
			if err != nil {
				os.Stderr.WriteString("Could not write licensor file.\n")
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}
	},
}
