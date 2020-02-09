package subcommands

import (
	"io"
	"licensezero.com/licensezero/user"
	"net/http"
)

const whoAmIDescription = "Show your identity information."

// WhoAmI prints identity information.
var WhoAmI = &Subcommand{
	Tag:         "misc",
	Description: whoAmIDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		identity, err := user.ReadIdentity()
		if err != nil {
			stderr.WriteString("Could not read identity file.\n")
			return 1
		}
		stdout.WriteString("Name: " + identity.Name + "\n")
		stdout.WriteString("Jurisdiction: " + identity.Jurisdiction + "\n")
		stdout.WriteString("E-Mail: " + identity.EMail + "\n")
		return 0
	},
}
