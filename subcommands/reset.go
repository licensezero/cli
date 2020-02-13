package subcommands

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
	"net/http"
)

const resetDescription = "Reset a seller access token."

var resetUsage = resetDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero reset [--broker URL]\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"broker URL": brokerFlagUsage,
	})

// Reset requests a seller token reset.
var Reset = &Subcommand{
	Tag:         "seller",
	Description: resetDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		// Parse flags.
		flagSet := flag.NewFlagSet("reset", flag.ContinueOnError)
		broker := brokerFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			stderr.WriteString(resetUsage)
		}
		err := flagSet.Parse(args)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				stderr.WriteString(err.Error() + "\n")
			}
			return 1
		}

		// Read saved accounts.
		accounts, err := user.ReadAccounts()
		if err != nil {
			stderr.WriteString("Error reading accounts: " + err.Error() + "\n")
			return 1
		}

		// Find account matching --broker.
		base := "https://" + *broker
		var brokerAccount *user.Account
		for _, account := range accounts {
			if account.Server == base {
				brokerAccount = account
				break
			}
		}
		if brokerAccount == nil {
			stderr.WriteString("No seller ID saved for " + *broker + "\n")
			return 1
		}
		sellerID := brokerAccount.SellerID
		stdout.WriteString("Broker Server: " + base + "\n")
		stdout.WriteString("Seller ID: " + sellerID + "\n")

		// Send request.
		brokerServer := api.BrokerServer{
			Client: client,
			Base:   base,
		}
		err = brokerServer.ResetToken(sellerID)
		if err != nil {
			stderr.WriteString("Error: " + err.Error() + "\n")
			return 1
		}
		stdout.WriteString("Check your e-mail for the reset link.\n")
		return 0
	},
}
