package subcommands

import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

const resetDescription = "Reset your licensor access token by e-mail."

var Reset = Subcommand{
	Description: resetDescription,
	Handler: func(args []string, paths Paths) {
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			os.Stderr.WriteString(identityHint + "\n")
			os.Exit(1)
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			os.Stderr.WriteString(licensorHint + "\n")
			os.Exit(1)
		}
		err = api.Reset(identity, licensor)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		os.Stdout.WriteString("Check your e-mail for the reset link.")
		os.Exit(0)
	},
}
