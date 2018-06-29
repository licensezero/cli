package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

const offerDescription = "Offer private licenses for sale."

var Offer = Subcommand{
	Tag:         "seller",
	Description: offerDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("offer", flag.ExitOnError)
		relicense := RelicenseFlag(flagSet)
		homepage := flagSet.String("homepage", "", "")
		description := flagSet.String("description", "", "")
		doNotOpen := DoNotOpen(flagSet)
		price := Price(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = offerUsage
		flagSet.Parse(args)
		if *price == 0 || *homepage == "" {
			offerUsage()
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		if !ConfirmAgencyTerms() {
			Fail(agencyTermsHint)
		}
		projectID, err := api.Offer(licensor, *homepage, *description, *price, *relicense)
		if err != nil {
			Fail(err.Error())
		}
		location := "https://licensezero.com/projects/" + projectID
		os.Stdout.WriteString("Project ID: " + projectID + "\n")
		openURLAndExit(location, doNotOpen)
	},
}

func offerUsage() {
	usage := offerDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero offer --price CENTS [--relicense CENTS] \\\n" +
		"                    --homepage URL --description TEXT\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"description TEXT": "Project description.",
			"do-not-open":      "Do not open project page in browser.",
			"homepage URL":     "Project homepage.",
			"price CENTS":      priceLine,
			"relicense CENTS":  relicenseLine,
		})
	Fail(usage)
}
