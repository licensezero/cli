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

const offerDescription = "Offer private licenses for sale."

const descriptionUsage = "Short project description."

const repositoryUsage = "Project source code URL."

var offerUsage = offerDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero offer --price CENTS [--relicense CENTS]\\\n" +
	"                    --repository URL --description TEXT\\\n" +
	"                    [--broker URL]\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"broker URL":       brokerFlagUsage,
		"description TEXT": descriptionUsage,
		"repository URL":   repositoryUsage,
		"price CENTS":      priceUsage,
		"relicense CENTS":  relicenseUsage,
	})

// Offer creates a project and offers private licenses for sale.
var Offer = &Subcommand{
	Tag:         "seller",
	Description: offerDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		// Parse flags.
		flagSet := flag.NewFlagSet("offer", flag.ExitOnError)
		broker := brokerFlag(flagSet)
		relicense := relicenseFlag(flagSet)
		repository := flagSet.String("repository", "", repositoryUsage)
		description := flagSet.String("description", "", descriptionUsage)
		price := priceFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			stderr.WriteString(offerUsage)
		}
		err := flagSet.Parse(args)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				stderr.WriteString(err.Error() + "\n")
			}
			return 1
		}
		if *price == 0 || *repository == "" {
			stderr.WriteString(offerUsage)
			return 1
		}

		// Find the relevant account.
		accounts, err := user.ReadAccounts()
		if err != nil {
			stderr.WriteString("Error reading accounts: " + err.Error() + "\n")
			return 1
		}
		if len(accounts) == 0 {
			stderr.WriteString("No accounts. Register for an account and run `licensezero token` first.\n")
			return 1
		}
		var account *user.Account
		if len(accounts) == 1 {
			account = accounts[0]
		}
		if len(accounts) > 1 {
			if *broker == "" {
				stderr.WriteString("You have multiple seller accounts.\nChoose one by passing --broker URL.\n")
				for _, account := range accounts {
					stderr.WriteString(account.Server + "\n")
				}
				return 1
			}
			base := "https://" + *broker
			for _, possible := range accounts {
				if possible.Server == base {
					account = possible
					break
				}
			}
			if account == nil {
				stderr.WriteString("No account for " + *broker + " found.\n")
				return 1
			}
		}

		// Agree to brokerage terms.
		confirmed, err := confirmBrokerageTerms(account.Server, stdin, stdout)
		if err != nil {
			return 1
		}
		if !confirmed {
			stderr.WriteString(brokerageTermsHint + "\n")
			return 1
		}

		// Send request.
		server := api.BrokerServer{
			Client: client,
			Base:   account.Server,
		}
		location, err := server.RegisterOffer(
			account.SellerID,
			account.Token,
			*repository,
			*description,
			*price,
			relicense,
		)
		if err != nil {
			stderr.WriteString("Error creating offer: " + err.Error() + "\n")
			return 1
		}
		stdout.WriteString("Created: " + location + "\n")
		return 0
	},
}
