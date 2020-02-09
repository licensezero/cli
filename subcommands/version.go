package subcommands

import (
	"io"
	"licensezero.com/licensezero/api"
)

const versionDescription = "Print version."

// Version prints the CLI version.
var Version = &Subcommand{
	Tag:         "misc",
	Description: versionDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client api.Client) int {
		if args[0] == "" {
			stdout.WriteString("Development Build\n")
		} else {
			stdout.WriteString(args[0])
		}
		return 0
	},
}
