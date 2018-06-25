package subcommands

import "fmt"
import "github.com/licensezero/cli/data"
import "os"

const whoAmIDescription = "Show your name, tax jurisdiction, and e-mail."

var WhoAmI = Subcommand{
	Tag:         "misc",
	Description: whoAmIDescription,
	Handler: func(args []string, paths Paths) {
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			os.Stderr.WriteString("Could not read identity file.\n")
			os.Exit(1)
		}
		licensor, err := data.ReadLicensor(paths.Home)
		fmt.Println("Name: " + identity.Name)
		fmt.Println("Jurisdiction: " + identity.Jurisdiction)
		fmt.Println("E-Mail: " + identity.EMail)
		if err == nil {
			fmt.Println("Licensor ID: " + licensor.LicensorID)
		}
		os.Exit(0)
	},
}
