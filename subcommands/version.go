package subcommands

import (
	"io"
)

const versionDescription = "Print version."

// Version prints the CLI version.
var Version = &Subcommand{
	Tag:         "misc",
	Description: versionDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter) int {
		if args[0] == "" {
			stdout.WriteString("Development Build\n")
		} else {
			stdout.WriteString(args[0])
		}
		return 0
	},
}
