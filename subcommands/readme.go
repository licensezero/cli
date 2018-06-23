package subcommands

import "os"

// TODO: Implement readme subcommand.

var README = Subcommand{
	Description: "Append licensing information to README.",
	Handler: func(args []string, paths Paths) {
		os.Exit(0)
	},
}
