package subcommands

import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const repriceDescription = "Change pricing."

// Reprice changes pricing.
var Reprice = &Subcommand{
	Description: repriceDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("reprice", flag.ExitOnError)
		price := priceFlag(flagSet)
		relicense := relicenseFlag(flagSet)
		noRelicense := noRelicenseFlag(flagSet)
		offerID := offerIDFlag(flagSet)
		id := idFlag(flagSet)
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = repriceUsage
		flagSet.Parse(args)
		if *price == 0 || (*offerID == "" && *id == "") {
			repriceUsage()
		}
		if *noRelicense && *relicense != 0 {
			repriceUsage()
		}
		if *offerID != "" && *id != "" {
			repriceUsage()
		}
		if *offerID != "" {
			*id = *offerID
		}
		if !validID(*id) {
			invalidID()
		}
		developer, err := data.ReadDeveloper(paths.Home)
		if err != nil {
			Fail(developerHint)
		}
		if err != nil {
			Fail(err.Error())
		}
		err = api.Reprice(developer, *id, *price, *relicense)
		if err != nil {
			Fail("Error sending reprice request:" + err.Error())
		}
		if !*silent {
			os.Stdout.WriteString("Repriced.\n")
		}
		os.Exit(0)
	},
}

func repriceUsage() {
	usage := repriceDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero reprice --id ID --price CENTS \\\n" +
		"                      (--relicense CENTS | --no-relicense)\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"price CENTS":     priceLine,
			"id ID":           idLine,
			"relicense CENTS": relicenseLine,
			"no-relicense":    noRelicenseLine,
			"silent":          silentLine,
		})
	Fail(usage)
}
