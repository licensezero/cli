package subcommands

import "os"

var Reset = Subcommand{
	Description: "Reset your access token.",
	Handler: func(args []string) {
		os.Exit(0)
	},
}
