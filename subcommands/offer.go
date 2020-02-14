package subcommands

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/api"
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
		account, message := selectAccount(broker)
		if message != "" {
			stderr.WriteString(message)
			return 1
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
