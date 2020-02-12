package subcommands

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/inventory"
	"licensezero.com/licensezero/user"
	"net/http"
	"os"
)

const quoteDescription = "Quote missing private licenses."

var quoteUsage = quoteDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero quote\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"json":          "Output JSON.",
		"noncommercial": noncommercialLine,
		"open":          openLine,
	})

// Quote generates a quote for missing private licenses.
var Quote = &Subcommand{
	Tag:         "buyer",
	Description: quoteDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		flagSet := flag.NewFlagSet("quote", flag.ExitOnError)
		noNoncommercial := noNoncommercialFlag(flagSet)
		noProsperity := noProsperityFlag(flagSet)
		noncommercial := noncommercialFlag(flagSet)
		noReciprocal := noReciprocalFlag(flagSet)
		open := openFlag(flagSet)
		noParity := noParityFlag(flagSet)
		outputJSON := flagSet.Bool("json", false, "")
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			stderr.WriteString(quoteUsage)
		}
		flagSet.Parse(args)
		ignoreNoncommercial := *noncommercial || *noNoncommercial || *noProsperity
		ignoreReciprocal := *open || *noReciprocal || *noParity

		// Compile inventory.
		configPath, err := user.ConfigPath()
		if err != nil {
			stderr.WriteString("Could not get configuration directory.\n")
			return 1
		}
		wd, err := os.Getwd()
		if err != nil {
			stderr.WriteString("Could not get working directory.\n")
			return 1
		}
		compiled, err := inventory.Compile(
			configPath,
			wd,
			ignoreNoncommercial,
			ignoreReciprocal,
			client,
		)
		if err != nil {
			stderr.WriteString("Error reading dependencies: " + err.Error() + "\n")
		}

		if *outputJSON {
			marshalled, err := json.Marshal(compiled)
			if err != nil {
				stderr.WriteString("Error serializing output: " + err.Error() + "\n")
				return 1
			}
			os.Stdout.WriteString(string(marshalled) + "\n")
			return 0
		}

		// Print summary.
		licensable := compiled.Licensable
		unlicensed := compiled.Unlicensed
		stdout.WriteString("Offers: " + string(len(licensable)) + "\n")
		stdout.WriteString("Licensed: " + string(len(compiled.Licensed)) + "\n")
		stdout.WriteString("Own: " + string(len(compiled.Own)) + "\n")
		stdout.WriteString("Unlicensed: " + string(len(unlicensed)) + "\n")
		stdout.WriteString("Ignored: " + string(len(compiled.Ignored)) + "\n")
		stdout.WriteString("Invalid: " + string(len(compiled.Invalid)) + "\n")
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
				Client: client,
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
			stdout.WriteString("\n")
			stdout.WriteString("- Server: " + result.Server + "\n")
			stdout.WriteString("  Offer: " + result.OfferID + "\n")
			stdout.WriteString("  URL: " + result.Offer.URL + "\n")
			seller := result.Seller
			if ok {
				stdout.WriteString("  Seller: " + seller.Name + " [" + seller.Jurisdiction + "] <" + seller.EMail + ">\n")
			}
			item := result.Item
			if item.Type != "" {
				stdout.WriteString("  Type: " + item.Type + "\n")
			}
			if item.Path != "" {
				stdout.WriteString("  Path: " + item.Path + "\n")
			}
			if item.Scope != "" {
				stdout.WriteString("  Scope: " + item.Scope + "\n")
			}
			if item.Name != "" {
				stdout.WriteString("  Name: " + item.Name + "\n")
			}
			if item.Version != "" {
				stdout.WriteString("  Version: " + item.Version + "\n")
			}
			stdout.WriteString("  Price: " + string(single.Amount) + " " + single.Currency)
		}

		stdout.WriteString("Totals:\n")
		for currency, total := range totals {
			stdout.WriteString("  " + string(total) + " " + currency + "\n")
		}

		if len(fetchErrors) != 0 {
			for _, err := range fetchErrors {
				stderr.WriteString("Error: " + err.Error() + "\n")
			}
			return 1
		}

		if len(unlicensed) > 0 {
			return 1
		}
		return 0
	},
}
