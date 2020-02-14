package subcommands

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/inventory"
	"licensezero.com/licensezero/user"
	"os"
	"strconv"
)

const quoteDescription = "Quote missing private licenses."

var quoteUsage = quoteDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero quote\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"json":          "Output JSON.",
		"noncommercial": noncommercialUsage,
		"open":          openUsage,
	})

// Quote generates a quote for missing private licenses.
var Quote = &Subcommand{
	Tag:         "buyer",
	Description: quoteDescription,
	Handler: func(env Environment) int {
		flagSet := flag.NewFlagSet("quote", flag.ExitOnError)
		noncommercial := noncommercialFlag(flagSet)
		open := openFlag(flagSet)
		outputJSON := flagSet.Bool("json", false, "")
		flagSet.SetOutput(ioutil.Discard)
		printUsage := func() {
			env.Stderr.WriteString(quoteUsage)
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

		if *outputJSON {
			marshalled, err := json.Marshal(compiled)
			if err != nil {
				env.Stderr.WriteString("Error serializing output: " + err.Error() + "\n")
				return 1
			}
			os.Stdout.WriteString(string(marshalled) + "\n")
			return 0
		}

		// Print summary.
		licensable := compiled.Licensable
		unlicensed := compiled.Unlicensed
		env.Stdout.WriteString("Offers: " + strconv.Itoa(len(licensable)) + "\n")
		env.Stdout.WriteString("Licensed: " + strconv.Itoa(len(compiled.Licensed)) + "\n")
		env.Stdout.WriteString("Own: " + strconv.Itoa(len(compiled.Own)) + "\n")
		env.Stdout.WriteString("Unlicensed: " + strconv.Itoa(len(unlicensed)) + "\n")
		env.Stdout.WriteString("Ignored: " + strconv.Itoa(len(compiled.Ignored)) + "\n")
		env.Stdout.WriteString("Invalid: " + strconv.Itoa(len(compiled.Invalid)) + "\n")
		if len(unlicensed) == 0 {
			return 0
		}

		// Calculate quote.
		type LineItem struct {
			Server  string
			OfferID string
			Public  string
			Item    inventory.Item
			Offer   *api.Offer
			Seller  *api.Seller
			Broker  *api.Broker
		}
		var results []LineItem
		var fetchErrors []error

		type SellerPointer struct {
			Server   string
			SellerID string
		}
		sellersCache := make(map[SellerPointer]*api.Seller)

		brokersCache := make(map[string]*api.Broker)

		for _, item := range unlicensed {
			server := api.BrokerServer{
				Client: env.Client,
				Base:   item.Server,
			}
			// Fetch offer.
			offer, err := server.Offer(item.OfferID)
			if err != nil {
				fetchErrors = append(fetchErrors, err)
				continue
			}
			// Fetch seller.
			seller, cached := sellersCache[SellerPointer{
				Server:   item.Server,
				SellerID: offer.SellerID,
			}]
			if !cached {
				seller, err = server.Seller(offer.SellerID)
				if err != nil {
					fetchErrors = append(fetchErrors, err)
				}
			}
			// Fetch broker.
			broker, cached := brokersCache[item.Server]
			if !cached {
				broker, err = server.Broker()
				if err != nil {
					fetchErrors = append(fetchErrors, err)
				}
			}
			results = append(results, LineItem{
				Server:  item.Server,
				OfferID: item.OfferID,
				Public:  item.Public,
				Item:    item,
				Offer:   offer,
				Seller:  seller,
				Broker:  broker,
			})
		}

		// Display line items and calculate totals.
		totals := make(map[string]uint)
		for _, result := range results {
			single := result.Offer.Pricing.Single
			total, ok := totals[single.Currency]
			if !ok {
				total = 0
			}
			total += single.Amount
			env.Stdout.WriteString("\n")
			env.Stdout.WriteString("- Server: " + result.Server + "\n")
			env.Stdout.WriteString("  Offer: " + result.OfferID + "\n")
			env.Stdout.WriteString("  URL: " + result.Offer.URL + "\n")
			seller := result.Seller
			if ok {
				env.Stdout.WriteString("  Seller: " + seller.Name + " [" + seller.Jurisdiction + "] <" + seller.EMail + ">\n")
			}
			item := result.Item
			if item.Type != "" {
				env.Stdout.WriteString("  Type: " + item.Type + "\n")
			}
			if item.Path != "" {
				env.Stdout.WriteString("  Path: " + item.Path + "\n")
			}
			if item.Scope != "" {
				env.Stdout.WriteString("  Scope: " + item.Scope + "\n")
			}
			if item.Name != "" {
				env.Stdout.WriteString("  Name: " + item.Name + "\n")
			}
			if item.Version != "" {
				env.Stdout.WriteString("  Version: " + item.Version + "\n")
			}
			env.Stdout.WriteString("  Price: " + string(single.Amount) + " " + single.Currency)
		}

		env.Stdout.WriteString("Totals:\n")
		for currency, total := range totals {
			env.Stdout.WriteString("  " + string(total) + " " + currency + "\n")
		}

		if len(fetchErrors) != 0 {
			for _, err := range fetchErrors {
				env.Stderr.WriteString("Error: " + err.Error() + "\n")
			}
			return 1
		}

		if len(unlicensed) > 0 {
			return 1
		}
		return 0
	},
}
