package subcommands

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"net/http"
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
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		flagSet := flag.NewFlagSet("reprice", flag.ExitOnError)
		broker := brokerFlag(flagSet)
		price := priceFlag(flagSet)
		relicense := relicenseFlag(flagSet)
		noRelicense := noRelicenseFlag(flagSet)
		offerID := offerIDFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			stderr.WriteString(repriceUsage)
		}
		err := flagSet.Parse(args)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				stderr.WriteString(err.Error() + "\n")
			}
			return 1
		}
		if *price == 0 || *offerID == "" {
			stderr.WriteString(repriceUsage)
			return 1
		}
		if *noRelicense && *relicense != 0 {
			stderr.WriteString(repriceUsage)
			return 1
		}
		if *offerID == "" {
			stderr.WriteString(repriceUsage)
			return 1
		}

		// Find the relevant account.
		account, message := selectAccount(broker)
		if message != "" {
			stderr.WriteString(message)
			return 1
		}

		server := api.BrokerServer{
			Client: client,
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
			stderr.WriteString("Error sending reprice request:" + err.Error() + "\n")
			return 1
		}
		stdout.WriteString("Repriced.\n")
		return 0
	},
}
