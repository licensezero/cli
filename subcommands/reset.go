package subcommands

import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "os"

const resetDescription = "Reset your API access token."

// Reset requests a new access token.
var Reset = &Subcommand{
	Tag:         "seller",
	Description: resetDescription,
	Handler: func(args []string, paths Paths) {
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			Fail(identityHint)
		}
		developer, err := data.ReadDeveloper(paths.Home)
		if err != nil {
			Fail(developerHint)
		}
		err = api.Reset(identity, developer)
		if err != nil {
			Fail("Error sending reset request: " + err.Error())
		}
		os.Stdout.WriteString("Check your e-mail for the reset link.\n")
		os.Exit(0)
	},
}
