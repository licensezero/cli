package subcommands

import "encoding/json"
import "flag"
import "fmt"
import "licensezero.com/cli/api"
import "licensezero.com/cli/inventory"
import "io/ioutil"
import "os"
import "strconv"

const quoteDescription = "Quote missing private licenses."

// Quote generates a quote for missing private licenses.
var Quote = &Subcommand{
	Tag:         "buyer",
	Description: quoteDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("quote", flag.ExitOnError)
		noNoncommercial := noNoncommercialFlag(flagSet)
		noProsperity := noProsperityFlag(flagSet)
		noncommercial := noncommercialFlag(flagSet)
		noReciprocal := noReciprocalFlag(flagSet)
		open := openFlag(flagSet)
		noParity := noParityFlag(flagSet)
		outputJSON := flagSet.Bool("json", false, "")
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = quoteUsage
		flagSet.Parse(args)
		suppressNoncommercial := *noncommercial || *noNoncommercial || *noProsperity
		suppressReciprocal := *open || *noReciprocal || *noParity
		offers, err := inventory.Inventory(paths.Home, paths.CWD, suppressNoncommercial, suppressReciprocal)
		if err != nil {
			Fail("Could not read dependeny tree.")
		}
		licensable := offers.Licensable
		licensed := offers.Licensed
		waived := offers.Waived
		unlicensed := offers.Unlicensed
		ignored := offers.Ignored
		invalid := offers.Invalid
		if *outputJSON {
			marshalled, err := json.Marshal(offers)
			if err != nil {
				Fail("Error serializing output.")
			}
			os.Stdout.WriteString(string(marshalled) + "\n")
			os.Exit(0)
		}
		if len(licensable) == 0 {
			fmt.Println("No License Zero dependencies found.")
			os.Exit(0)
		}
		fmt.Printf("License Zero Offers: %d\n", len(licensable))
		fmt.Printf("Licensed: %d\n", len(licensed))
		fmt.Printf("Waived: %d\n", len(waived))
		fmt.Printf("Ignored: %d\n", len(ignored))
		fmt.Printf("Unlicensed: %d\n", len(unlicensed))
		fmt.Printf("Invalid: %d\n", len(invalid))
		if len(unlicensed) == 0 {
			os.Exit(0)
		}
		var offerIDs []string
		for _, offer := range unlicensed {
			offerIDs = append(offerIDs, offer.OfferID)
		}
		response, err := api.Quote(offerIDs)
		if err != nil {
			Fail("Error requesting quote: " + err.Error())
		}
		var total uint
		for _, offer := range response {
			var prior *inventory.Offer
			for _, candidate := range unlicensed {
				if candidate.OfferID == offer.OfferID {
					prior = &candidate
					break
				}
			}
			total += offer.Pricing.Private
			fmt.Println("\n- Offer: " + offer.OfferID)
			fmt.Println("  Description: " + offer.Description)
			fmt.Println("  URL: " + offer.URL)
			if prior != nil {
				var terms = prior.License.Terms
				if terms == "noncommercial" {
					fmt.Println("  Terms: Noncommercial")
				} else if terms == "reciprocal" {
					fmt.Println("  Terms: Reciprocal")
				} else if terms == "parity" {
					fmt.Println("  Terms: Parity")
				} else if terms == "prosperity" {
					fmt.Println("  Terms: Prosperity")
				}
			}
			fmt.Println("  Licensor: " + offer.Licensor.Name + " [" + offer.Licensor.Jurisdiction + "]")
			if offer.Retracted {
				fmt.Println("  Retracted!")
			}
			if prior != nil {
				var artifact = prior.Artifact
				if artifact.Type != "" {
					fmt.Println("  Type: " + artifact.Type)
				}
				if artifact.Path != "" {
					fmt.Println("  Path: " + artifact.Path)
				}
				if artifact.Scope != "" {
					fmt.Println("  Scope: " + artifact.Scope)
				}
				if artifact.Name != "" {
					fmt.Println("  Name: " + artifact.Name)
				}
				if artifact.Version != "" {
					fmt.Println("  Version: " + artifact.Version)
				}
			}
			fmt.Println("  Price: " + currency(offer.Pricing.Private))
			fmt.Printf("\nTotal: %s\n", currency(total))
		}
		if len(unlicensed) == 0 {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	},
}

func quoteUsage() {
	usage := quoteDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero quote\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"json":          "Output JSON.",
			"noncommercial": noncommercialLine,
			"open":          openLine,
		})
	Fail(usage)
}

func currency(cents uint) string {
	if cents < 100 {
		if cents < 10 {
			return "$0.0" + strconv.Itoa(int(cents))
		}
		return "$0." + strconv.Itoa(int(cents))
	}
	asString := fmt.Sprintf("%d", cents)
	return "$" + asString[:len(asString)-2] + "." + asString[len(asString)-2:]
}
