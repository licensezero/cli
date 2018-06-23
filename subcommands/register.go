package subcommands

import "os"

// TODO: Implement register subcommand.

var Register = Subcommand{
	Description: "Register to sell private licenses through licensezero.com.",
	Handler: func(args []string, paths Paths) {
		os.Exit(0)
	},
}
