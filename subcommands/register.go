package subcommands

import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

const registerDescription = "Register to sell private licenses through licensezero.com."

var Register = Subcommand{
	Description: registerDescription,
	Handler: func(args []string, paths Paths) {
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			os.Stderr.WriteString(identityHint)
			os.Exit(1)
		}
		os.Stdout.WriteString("Name: " + identity.Name + "\n")
		os.Stdout.WriteString("Jurisdiction: " + identity.Jurisdiction + "\n")
		os.Stdout.WriteString("E-Mail: " + identity.EMail + "\n")
		if !Confirm("Is this information correct?") {
			os.Stdout.WriteString("Exiting.\n")
			os.Exit(1)
		}
		if !ConfirmTermsOfService() {
			os.Stderr.WriteString("You must agree to the terms of service to register.\n")
			os.Exit(1)
		}
		err = api.Register(identity)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		os.Stdout.WriteString("Follow the Stripe authorization link sent by e-mail.\n")
		os.Exit(0)
	},
}

func registerUsage() {
	usage := registerDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero register\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
