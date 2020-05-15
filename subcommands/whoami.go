package subcommands

import "fmt"
import "licensezero.com/cli/data"
import "os"

const whoAmIDescription = "Show your identity information."

// WhoAmI prints identity information.
var WhoAmI = &Subcommand{
	Tag:         "misc",
	Description: whoAmIDescription,
	Handler: func(args []string, paths Paths) {
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			Fail("Could not read identity file.")
		}
		developer, err := data.ReadDeveloper(paths.Home)
		fmt.Println("Name: " + identity.Name)
		fmt.Println("Jurisdiction: " + identity.Jurisdiction)
		fmt.Println("E-Mail: " + identity.EMail)
		if err == nil {
			fmt.Println("Developer ID: " + developer.DeveloperID)
		}
		os.Exit(0)
	},
}
