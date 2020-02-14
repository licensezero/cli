package subcommands

import (
	"errors"
	"flag"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
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
	Handler: func(env Environment) int {
		// Parse flags.
		flagSet := flag.NewFlagSet("reset", flag.ContinueOnError)
		broker := brokerFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		printUsage := func() {
			env.Stderr.WriteString(resetUsage)
		}
		flagSet.Usage = printUsage
		err := flagSet.Parse(env.Arguments)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				printUsage()
			}
			return 1
		}

		// Read saved accounts.
		accounts, err := user.ReadAccounts()
		if err != nil {
			env.Stderr.WriteString("Error reading accounts: " + err.Error() + "\n")
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
			env.Stderr.WriteString("No seller ID saved for " + *broker + "\n")
			return 1
		}
		sellerID := brokerAccount.SellerID
		env.Stdout.WriteString("Broker Server: " + base + "\n")
		env.Stdout.WriteString("Seller ID: " + sellerID + "\n")

		// Send request.
		brokerServer := api.BrokerServer{
			Client: env.Client,
			Base:   base,
		}
		err = brokerServer.ResetToken(sellerID)
		if err != nil {
			env.Stderr.WriteString("Error: " + err.Error() + "\n")
			return 1
		}
		env.Stdout.WriteString("Check your e-mail for the reset link.\n")
		return 0
	},
}
