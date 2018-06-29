package subcommands

import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

const registerDescription = "Register to sell private licenses."

var Register = Subcommand{
	Tag:         "seller",
	Description: registerDescription,
	Handler: func(args []string, paths Paths) {
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			Fail(identityHint)
		}
		os.Stdout.WriteString("Name: " + identity.Name + "\n")
		os.Stdout.WriteString("Jurisdiction: " + identity.Jurisdiction + "\n")
		os.Stdout.WriteString("E-Mail: " + identity.EMail + "\n")
		if !Confirm("Is this information correct?") {
			os.Stdout.WriteString("Exiting.\n")
			os.Exit(1)
		}
		if !ConfirmTermsOfService() {
			Fail(termsHint)
		}
		err = api.Register(identity)
		if err != nil {
			Fail(err.Error())
		}
		os.Stdout.WriteString("Follow the Stripe authorization link sent by e-mail.\n")
		os.Exit(0)
	},
}

func registerUsage() {
	usage := registerDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero register\n"
	Fail(usage)
}
