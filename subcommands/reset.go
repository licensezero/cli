package subcommands

import "os"

var Reset = Subcommand{
	Description: "Reset your access token.",
	Handler: func(args []string, home string) {
		os.Exit(0)
	},
}
