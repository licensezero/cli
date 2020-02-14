package subcommands

import (
	"errors"
	"flag"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/inventory"
	"licensezero.com/licensezero/user"
	"os"
)

const buyDescription = "Buy missing private licenses."

var buyUsage = buyDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero buy\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"noncommercial": noncommercialUsage,
		"open":          openUsage,
		"do-not-open":   doNotOpenUsage,
	})

// Buy opens a buy page for each broker.
var Buy = &Subcommand{
	Tag:         "buyer",
	Description: buyDescription,
	Handler: func(env Environment) int {
		flagSet := flag.NewFlagSet("buy", flag.ExitOnError)
		doNotOpen := doNotOpenFlag(flagSet)
		noncommercial := noncommercialFlag(flagSet)
		open := openFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		printUsage := func() {
			env.Stderr.WriteString(buyUsage)
		}
		flagSet.Usage = printUsage
		err := flagSet.Parse(env.Arguments)
		if err != nil {
			if errors.Is(err, flag.ErrHelp) {
				printUsage()
			}
			return 1
		}

		// Compile inventory.
		configPath, err := user.ConfigPath()
		if err != nil {
			env.Stderr.WriteString("Could not get configuration directory.\n")
			return 1
		}
		wd, err := os.Getwd()
		if err != nil {
			env.Stderr.WriteString("Could not get working directory.\n")
			return 1
		}
		compiled, err := inventory.Compile(
			configPath,
			wd,
			*noncommercial,
			*open,
			env.Client,
		)
		if err != nil {
			env.Stderr.WriteString("Error reading dependencies: " + err.Error() + "\n")
			return 1
		}

		licensable := compiled.Licensable
		unlicensed := compiled.Unlicensed
		if len(licensable) == 0 {
			env.Stdout.WriteString("No License Zero artifacts found.\n")
			return 0
		}
		if len(unlicensed) == 0 {
			env.Stdout.WriteString("No private licenses to buy.\n")
			return 0
		}

		// Create a map from broker server URL to slice of offerIDs,
		// so that we can create one order per broker server for all
		// licenses needed from that server.
		servers := make(map[string][]string)
		for _, item := range unlicensed {
			server := item.Server
			offerIDs, ok := servers[server]
			if !ok {
				servers[server] = []string{}
			}
			servers[server] = append(offerIDs, item.OfferID)
		}

		// Send order requests to broker servers.
		identity, err := user.ReadIdentity()
		if err != nil {
			env.Stderr.WriteString(identityHint)
			return 1
		}
		hadOrderError := false
		for base, offerIDs := range servers {
			server := api.BrokerServer{
				Client: env.Client,
				Base:   base,
			}
			location, err := server.Order(
				identity.EMail,
				identity.Jurisdiction,
				identity.Name,
				offerIDs,
			)
			if err != nil {
				hadOrderError = true
				env.Stderr.WriteString("Error ordering from " + base + ": " + err.Error() + "\n")
				continue
			}
			openURL(location, doNotOpen, env.Stdout)
		}
		if hadOrderError {
			return 1
		}
		return 0
	},
}
