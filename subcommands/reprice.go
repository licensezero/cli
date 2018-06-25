package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

const repriceDescription = "Change pricing for your project."

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
			os.Stderr.WriteString(licensorHint + "\n")
			os.Exit(1)
		}
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
		var projectID string
		if *projectIDFlag != "" {
			projectID = *projectIDFlag
		} else {
			projectIDs, _, err := readEntries(paths.CWD)
			if err != nil {
				os.Stderr.WriteString("Could not read package.json.\n")
				os.Exit(1)
			}
			if len(projectIDs) > 0 {
				os.Stderr.WriteString("package.json has metadata for multiple License Zero projects.\n")
				os.Stderr.WriteString("Use --project to specify your project ID.")
				os.Exit(1)
			}
		}
		err = api.Reprice(licensor, projectID, *price, *relicense)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
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
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
