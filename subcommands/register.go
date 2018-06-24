package subcommands

import "os"

// TODO: Implement register subcommand.

const registerDescription = "Register to sell private licenses through licensezero.com"

var Register = Subcommand{
	Description: registerDescription,
	Handler: func(args []string, paths Paths) {
		os.Exit(0)
	},
}
