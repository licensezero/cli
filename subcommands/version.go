package subcommands

import (
	"os"
)

const versionDescription = "Print version."

// Version prints the CLI version.
var Version = &Subcommand{
	Tag:         "misc",
	Description: versionDescription,
	Handler: func(args []string, stdin, stdout, stderr *os.File) int {
		if args[0] == "" {
			stdout.WriteString("Development Build\n")
		} else {
			stdout.WriteString(args[0])
		}
		return 0
	},
}
