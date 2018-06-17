package subcommands

import "os"

var WhoAmI = Subcommand{
	Description: "Show your licensor ID.",
	Handler: func(args []string) {
		os.Exit(0)
	},
}
