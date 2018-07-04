package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
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
		projectIDFlag := projectIDFlag(flagSet)
		justIDFlag := idFlag(flagSet)
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = repriceUsage
		flagSet.Parse(args)
		if *price == 0 || (*projectIDFlag == "" && *justIDFlag == "") {
			repriceUsage()
		}
		if *projectIDFlag != "" && *justIDFlag != "" {
			repriceUsage()
		}
		if *projectIDFlag != "" {
			*justIDFlag = *projectIDFlag
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		if err != nil {
			Fail(err.Error())
		}
		var id string
		if *justIDFlag != "" {
			id = *projectIDFlag
		} else {
			projectIDs, _, err := readEntries(paths.CWD)
			if err != nil {
				Fail("Could not read licensezero.json.")
			}
			if len(projectIDs) > 0 {
				os.Stderr.WriteString("licensezero.json has metadata for multiple License Zero projects.\n")
				Fail("Use --id to specify your project ID.")
			}
		}
		err = api.Reprice(licensor, id, *price, *relicense)
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
		"  licensezero reprice --id ID --price CENTS [--relicense CENTS]\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"price CENTS":     priceLine,
			"id ID":           idLine,
			"relicense CENTS": relicenseLine,
			"silent":          silentLine,
		})
	Fail(usage)
}
