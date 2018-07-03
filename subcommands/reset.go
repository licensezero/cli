package subcommands

import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
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
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		err = api.Reset(identity, licensor)
		if err != nil {
			Fail("Error sending reset request: " + err.Error())
		}
		os.Stdout.WriteString("Check your e-mail for the reset link.\n")
		os.Exit(0)
	},
}
