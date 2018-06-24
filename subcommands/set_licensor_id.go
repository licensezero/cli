package subcommands

import "flag"
import "fmt"
import "github.com/licensezero/cli/data"
import "os"

const setLicensorIDDescription = "Set your licensezero.com licensor ID"

var SetLicensorID = Subcommand{
	Description: setLicensorIDDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("quote", flag.ContinueOnError)
		licensorID := flag.String("licensor-id", "", "")
		err := flagSet.Parse(args)
		if err != nil || *licensorID == "" {
			setLicensorIDUsage()
		}
		if len(args) != 1 {
			setLicensorIDUsage()
		}
		fmt.Println(licensorID)
		token := SecretPrompt("Token: ")
		newLicensor := data.Licensor{
			LicensorID: *licensorID,
			Token:      token,
		}
		existingLicensor, _ := data.ReadLicensor(paths.Home)
		if existingLicensor != nil && *existingLicensor != newLicensor {
			if !Confirm("Overwrite existing licensor info?") {
				os.Exit(0)
			}
		}
		err = data.WriteLicensor(paths.Home, &newLicensor)
		if err != nil {
			os.Stderr.WriteString("Could not write licensor file.\n")
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	},
}

func setLicensorIDUsage() {
	usage := setLicensorIDDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero set-licensor-id --licensor-id ID\n\n" +
		"Options:\n" +
		"  --licensor-id ID  Licensor ID (UUID)."
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
