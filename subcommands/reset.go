package subcommands

import "os"

// TODO: Implement reset subcommand.

var Reset = Subcommand{
	Description: "Reset your access token.",
	Handler: func(args []string, paths Paths) {
		os.Exit(0)
	},
}
