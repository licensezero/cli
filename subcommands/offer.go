package subcommands

import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const offerDescription = "Offer private licenses for sale."

// Offer creates a project and offers private licenses for sale.
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
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		if !confirmAgencyTerms() {
			Fail(agencyTermsHint)
		}
		projectID, err := api.Offer(licensor, *repository, *description, *price, *relicense)
		if err != nil {
			Fail("Error sending offer request: " + err.Error())
		}
		location := "https://licensezero.com/projects/" + projectID
		os.Stdout.WriteString("Project ID: " + projectID + "\n")
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
