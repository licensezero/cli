package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

const repriceDescription = "Change project pricing."

// TODO: Clarify UI for withdrawing relicense offers.

var Reprice = Subcommand{
	Tag:         "seller",
	Description: repriceDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("reprice", flag.ExitOnError)
		price := Price(flagSet)
		relicense := Relicense(flagSet)
		projectIDFlag := ProjectID(flagSet)
		silent := Silent(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = repriceUsage
		flagSet.Parse(args)
		if *price == 0 || *projectIDFlag == "" {
			repriceUsage()
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		if err != nil {
			Fail(err.Error())
		}
		var projectID string
		if *projectIDFlag != "" {
			projectID = *projectIDFlag
		} else {
			projectIDs, _, err := readEntries(paths.CWD)
			if err != nil {
				Fail("Could not read package.json.")
			}
			if len(projectIDs) > 0 {
				os.Stderr.WriteString("package.json has metadata for multiple License Zero projects.\n")
				Fail("Use --project to specify your project ID.")
			}
		}
		err = api.Reprice(licensor, projectID, *price, *relicense)
		if err != nil {
			Fail(err.Error())
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
		"  licensezero reprice --project ID --price CENTS [--relicense CENTS]\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"price":      priceLine,
			"project ID": "Project ID (UUID).",
			"relicense":  relicenseLine,
			"silent":     silentLine,
		})
	Fail(usage)
}
