package subcommands

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"net/http"
)

const lockDescription = "Lock pricing and availability."

const unlockUsage = "Unlock date and time (RFC 3339/ISO 8601)."

var lockUsage = lockDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero lock --offer ID --unlock DATE\n\n" +
	"                   [--broker URL]\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"broker URL":      brokerFlagUsage,
		"offer ID":        offerIDUsage,
		"unlock DATETIME": unlockUsage,
	})

// Lock fixes pricing and availability.
var Lock = &Subcommand{
	Tag:         "seller",
	Description: lockDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		flagSet := flag.NewFlagSet("lock", flag.ExitOnError)
		broker := brokerFlag(flagSet)
		offerID := offerIDFlag(flagSet)
		unlock := flagSet.String("unlock", "", unlockUsage)
		flagSet.SetOutput(ioutil.Discard)
		printUsage := func() {
			stderr.WriteString(lockUsage)
		}
		flagSet.Usage = printUsage
		err := flagSet.Parse(args)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				stderr.WriteString(err.Error() + "\n")
			}
			return 1
		}
		if *unlock == "" || *offerID == "" {
			printUsage()
			return 1
		}

		account, message := selectAccount(broker)
		if message != "" {
			stderr.WriteString(message)
			return 1
		}

		server := api.BrokerServer{
			Client: client,
			Base:   account.Server,
		}
		err = server.Lock(
			account.SellerID,
			account.Token,
			*offerID,
			*unlock,
		)
		if err != nil {
			stderr.WriteString("Error sending lock request: " + err.Error() + "\n")
			return 1
		}
		stdout.WriteString("Locked pricing.\n")
		return 0
	},
}
