package subcommands

import (
	"fmt"
	"licensezero.com/licensezero/user"
	"os"
)

const whoAmIDescription = "Show your identity information."

// WhoAmI prints identity information.
var WhoAmI = &Subcommand{
	Tag:         "misc",
	Description: whoAmIDescription,
	Handler: func(args []string, paths Paths) {
		identity, err := user.ReadIdentity(paths.Home)
		if err != nil {
			Fail("Could not read identity file.")
		}
		fmt.Println("Name: " + identity.Name)
		fmt.Println("Jurisdiction: " + identity.Jurisdiction)
		fmt.Println("E-Mail: " + identity.EMail)
		os.Exit(0)
	},
}
