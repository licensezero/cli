package subcommands

import "os"

var Reset = Subcommand{
	Description: "Reset your access token.",
	Handler: func(args []string, paths Paths) {
		os.Exit(0)
	},
}
