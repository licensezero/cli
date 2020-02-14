package subcommands

import (
	"errors"
	"flag"
	"io/ioutil"
	"licensezero.com/licensezero/api"
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
	Handler: func(env Environment) int {
		// Parse flags.
		flagSet := flag.NewFlagSet("offer", flag.ExitOnError)
		broker := brokerFlag(flagSet)
		relicense := relicenseFlag(flagSet)
		repository := flagSet.String("repository", "", repositoryUsage)
		description := flagSet.String("description", "", descriptionUsage)
		price := priceFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			env.Stderr.WriteString(offerUsage)
		}
		err := flagSet.Parse(env.Arguments)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				env.Stderr.WriteString(err.Error() + "\n")
			}
			return 1
		}
		if *price == 0 || *repository == "" {
			env.Stderr.WriteString(offerUsage)
			return 1
		}

		// Find the relevant account.
		account, message := selectAccount(broker)
		if message != "" {
			env.Stderr.WriteString(message)
			return 1
		}

		// Agree to brokerage terms.
		confirmed, err := confirmBrokerageTerms(account.Server, env.Stdin, env.Stdout)
		if err != nil {
			return 1
		}
		if !confirmed {
			env.Stderr.WriteString(brokerageTermsHint + "\n")
			return 1
		}

		// Send request.
		server := api.BrokerServer{
			Client: env.Client,
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
			env.Stderr.WriteString("Error creating offer: " + err.Error() + "\n")
			return 1
		}
		env.Stdout.WriteString("Created: " + location + "\n")
		return 0
	},
}
