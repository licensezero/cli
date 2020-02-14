package subcommands

import (
	"errors"
	"flag"
	"io/ioutil"
	"licensezero.com/licensezero/api"
)

const repriceDescription = "Change pricing."

var repriceUsage = repriceDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero reprice --offer UUID --price CENTS \\\n" +
	"                      (--relicense CENTS | --no-relicense)\\\n" +
	"                      [--broker URL]\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"broker URL":      brokerFlagUsage,
		"price CENTS":     priceUsage,
		"offer UUID":      offerIDUsage,
		"relicense CENTS": relicenseUsage,
		"no-relicense":    noRelicenseUsage,
	})

// Reprice changes pricing.
var Reprice = &Subcommand{
	Tag:         "seller",
	Description: repriceDescription,
	Handler: func(env Environment) int {
		flagSet := flag.NewFlagSet("reprice", flag.ExitOnError)
		broker := brokerFlag(flagSet)
		price := priceFlag(flagSet)
		relicense := relicenseFlag(flagSet)
		noRelicense := noRelicenseFlag(flagSet)
		offerID := offerIDFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			env.Stderr.WriteString(repriceUsage)
		}
		err := flagSet.Parse(env.Arguments)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				env.Stderr.WriteString(err.Error() + "\n")
			}
			return 1
		}
		if *price == 0 || *offerID == "" {
			env.Stderr.WriteString(repriceUsage)
			return 1
		}
		if *noRelicense && *relicense != 0 {
			env.Stderr.WriteString(repriceUsage)
			return 1
		}
		if *offerID == "" {
			env.Stderr.WriteString(repriceUsage)
			return 1
		}

		// Find the relevant account.
		account, message := selectAccount(broker)
		if message != "" {
			env.Stderr.WriteString(message)
			return 1
		}

		server := api.BrokerServer{
			Client: env.Client,
			Base:   account.Server,
		}
		err = server.Reprice(
			account.SellerID,
			account.Token,
			*offerID,
			*price,
			relicense,
		)
		if err != nil {
			env.Stderr.WriteString("Error sending reprice request:" + err.Error() + "\n")
			return 1
		}
		env.Stdout.WriteString("Repriced.\n")
		return 0
	},
}
