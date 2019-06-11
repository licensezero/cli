package subcommands

import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const repriceDescription = "Change project pricing."

// Reprice changes project pricing.
var Reprice = &Subcommand{
	Tag:         "seller",
	Description: repriceDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("reprice", flag.ExitOnError)
		price := priceFlag(flagSet)
		relicense := relicenseFlag(flagSet)
		noRelicense := noRelicenseFlag(flagSet)
		projectID := projectIDFlag(flagSet)
		id := idFlag(flagSet)
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = repriceUsage
		flagSet.Parse(args)
		if *price == 0 || (*projectID == "" && *id == "") {
			repriceUsage()
		}
		if *noRelicense && *relicense != 0 {
			repriceUsage()
		}
		if *projectID != "" && *id != "" {
			repriceUsage()
		}
		if *projectID != "" {
			*id = *projectID
		}
		if !validID(*id) {
			invalidID()
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		if err != nil {
			Fail(err.Error())
		}
		err = api.Reprice(licensor, *id, *price, *relicense)
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
