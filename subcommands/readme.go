package subcommands

import "os"

var README = Subcommand{
	Description: "Append licensing information to README.",
	Handler: func(args []string, home string) {
		os.Exit(0)
	},
}
