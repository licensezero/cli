package subcommands

import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const offerDescription = "Offer private licenses for sale."

// Offer creates an offer and offers private licenses for sale.
var Offer = &Subcommand{
	Tag:         "seller",
	Description: offerDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("offer", flag.ExitOnError)
		relicense := relicenseFlag(flagSet)
		noRelicense := flagSet.Bool("no-relicense", false, "")
		repository := flagSet.String("repository", "", "")
		description := flagSet.String("description", "", "")
		doNotOpen := doNotOpenFlag(flagSet)
		price := priceFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = offerUsage
		flagSet.Parse(args)
		if *price == 0 || *repository == "" {
			offerUsage()
		}
		if *noRelicense && *relicense != 0 {
			offerUsage()
		}
		developer, err := data.ReadDeveloper(paths.Home)
		if err != nil {
			Fail(developerHint)
		}
		if !confirmAgencyTerms() {
			Fail(agencyTermsHint)
		}
		offerID, err := api.Offer(developer, *repository, *description, *price, *relicense)
		if err != nil {
			Fail("Error sending offer request: " + err.Error())
		}
		location := "https://licensezero.com/offers/" + offerID
		os.Stdout.WriteString("Offer ID: " + offerID + "\n")
		openURLAndExit(location, doNotOpen)
	},
}

func offerUsage() {
	usage := offerDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero offer --price CENTS (--relicense CENTS || --no-relicense)\\\n" +
		"                    --repository URL --description TEXT\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"description TEXT": "Description.",
			"do-not-open":      "Do not open page in browser.",
			"repository URL":   "Source code repository URL.",
			"price CENTS":      priceLine,
			"relicense CENTS":  relicenseLine,
			"no-relicense":     noRelicenseLine,
		})
	Fail(usage)
}
