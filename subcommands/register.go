package subcommands

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
	"net/http"
	"os"
)

const registerDescription = "Register to sell private licenses."

const defaultBroker = "broker.licensezero.com"

var brokerFlagUsage = "Broker server name [default: " + defaultBroker + "]"

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
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		flagSet := flag.NewFlagSet("register", flag.ContinueOnError)
		broker := flagSet.String("broker", defaultBroker, brokerFlagUsage)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			stderr.WriteString(registerUsage)
		}
		err := flagSet.Parse(args)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				stderr.WriteString("\nError: " + err.Error() + "\n")
			}
			return 1
		}
		base := "https://" + *broker
		identity, err := user.ReadIdentity()
		if err != nil {
			stderr.WriteString(identityHint)
			return 1
		}
		os.Stdout.WriteString("Name: " + identity.Name + "\n")
		os.Stdout.WriteString("Jurisdiction: " + identity.Jurisdiction + "\n")
		os.Stdout.WriteString("E-Mail: " + identity.EMail + "\n")
		confirmed, err := stdin.Confirm("Is this information correct?", stdout)
		if err != nil {
			stderr.WriteString(err.Error())
			return 1
		}
		if !confirmed {
			os.Stdout.WriteString("Exiting.\n")
			return 1
		}
		confirmed, err = confirmTermsOfService(base, stdin, stdout)
		if err != nil {
			stderr.WriteString(err.Error())
			return 1
		}
		if !confirmed {
			os.Stdout.WriteString(termsHint)
			return 1
		}
		server := api.BrokerServer{
			Client: client,
			Base:   base,
		}
		err = server.RegisterSeller(
			identity.EMail,
			identity.Jurisdiction,
			identity.Name,
		)
		if err != nil {
			stderr.WriteString("Error sending register request: " + err.Error())
			return 1
		}
		stdout.WriteString("Follow the authorization link sent by e-mail.\n")
		stdout.WriteString("If you cannot find the e-mail, check your junk mail folder.\n")
		return 0
	},
}
