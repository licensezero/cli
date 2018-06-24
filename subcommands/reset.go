package subcommands

import "os"

// TODO: Implement reset subcommand.

const resetDescription = "Reset your access token."

var Reset = Subcommand{
	Description: resetDescription,
	Handler: func(args []string, paths Paths) {
		os.Exit(0)
	},
}
