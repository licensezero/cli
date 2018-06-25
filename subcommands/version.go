package subcommands

import "os"

const versionDescription = "Print version."

var Version = Subcommand{
	Tag:         "misc",
	Description: versionDescription,
	Handler: func(args []string, paths Paths) {
		if args[0] == "" {
			os.Stdout.WriteString("Development Version\n")
		} else {
			os.Stdout.WriteString(args[0] + "\n")
		}
		os.Exit(0)
	},
}
