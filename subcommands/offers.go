package subcommands

import "encoding/json"
import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const offersDescription = "List your offers."

// Offers prints the licensor's offers.
var Offers = &Subcommand{
	Tag:         "misc",
	Description: offersDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("offers", flag.ExitOnError)
		retracted := flagSet.Bool("include-retracted", false, "")
		outputJSON := flagSet.Bool("json", false, "")
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = offersUsage
		flagSet.Parse(args)
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		_, offers, err := api.Licensor(licensor.LicensorID)
		if err != nil {
			Fail("Could not fetch licensor information: " + err.Error())
		}
		var filtered []api.OfferInformation
		if *retracted {
			filtered = offers
		} else {
			for _, offer := range offers {
				if offer.Retracted == "" {
					filtered = append(filtered, offer)
				}
			}
		}
		if *outputJSON {
			marshalled, err := json.Marshal(filtered)
			if err != nil {
				Fail("Error serializing output.")
			}
			os.Stdout.WriteString(string(marshalled) + "\n")
			os.Exit(0)
		}
		for i, offer := range filtered {
			if i != 0 {
				os.Stdout.WriteString("\n")
			}
			os.Stdout.WriteString("- Offer ID: " + offer.OfferID + "\n")
			os.Stdout.WriteString("  Offered:    " + offer.Offered + "\n")
			if offer.Retracted != "" {
				os.Stdout.WriteString("  Retracted:  " + offer.Offered + "\n")
			}
		}
		os.Exit(0)
	},
}

func offersUsage() {
	usage := offersDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero offers\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"json":              "Output JSON.",
			"include-retracted": "List retracted offers.",
		})
	Fail(usage)
}
