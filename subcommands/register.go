package subcommands

import (
	"errors"
	"flag"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
)

const registerDescription = "Register to sell private licenses."

var registerUsage = registerDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero register --broker SERVER\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"broker URL": brokerFlagUsage,
	})

// Register a user to sell private licenses.
var Register = &Subcommand{
	Tag:         "seller",
	Description: registerDescription,
	Handler: func(env Environment) int {
		flagSet := flag.NewFlagSet("register", flag.ContinueOnError)
		broker := brokerFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			env.Stderr.WriteString(registerUsage)
		}
		err := flagSet.Parse(env.Arguments)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				env.Stderr.WriteString("\nError: " + err.Error() + "\n")
			}
			return 1
		}
		base := "https://" + *broker
		identity, err := user.ReadIdentity()
		if err != nil {
			env.Stderr.WriteString(identityHint)
			return 1
		}
		env.Stdout.WriteString("Name: " + identity.Name + "\n")
		env.Stdout.WriteString("Jurisdiction: " + identity.Jurisdiction + "\n")
		env.Stdout.WriteString("E-Mail: " + identity.EMail + "\n")
		confirmed, err := env.Stdin.Confirm("Is this information correct?", env.Stdout)
		if err != nil {
			env.Stderr.WriteString(err.Error())
			return 1
		}
		if !confirmed {
			env.Stdout.WriteString("Exiting.\n")
			return 1
		}
		confirmed, err = confirmTermsOfService(base, env.Stdin, env.Stdout)
		if err != nil {
			env.Stderr.WriteString(err.Error())
			return 1
		}
		if !confirmed {
			env.Stdout.WriteString(termsHint)
			return 1
		}
		server := api.BrokerServer{
			Client: env.Client,
			Base:   base,
		}
		err = server.RegisterSeller(
			identity.EMail,
			identity.Jurisdiction,
			identity.Name,
		)
		if err != nil {
			env.Stderr.WriteString("Error sending register request: " + err.Error())
			return 1
		}
		env.Stdout.WriteString("Follow the authorization link sent by e-mail.\n")
		env.Stdout.WriteString("If you cannot find the e-mail, check your junk mail folder.\n")
		return 0
	},
}
