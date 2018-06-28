package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

const waiveDescription = "Generate a waiver for your project."

var Waive = Subcommand{
	Tag:         "seller",
	Description: waiveDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("waive", flag.ExitOnError)
		jurisdiction := flagSet.String("jurisdiction", "", "Jurisdiction.")
		days := flagSet.Uint("days", 0, "Days.")
		forever := flagSet.Bool("forever", false, "Forever.")
		beneficiary := flagSet.String("beneficiary", "", "Beneficiary legal name.")
		projectID := ProjectID(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = waiveUsage
		flagSet.Parse(args)
		if *projectID == "" {
			waiveUsage()
		} else if *forever && *days > 0 {
			waiveUsage()
		} else if *days == 0 && !*forever {
			waiveUsage()
		} else if *beneficiary == "" || *jurisdiction == "" {
			waiveUsage()
		} else {
			licensor, err := data.ReadLicensor(paths.Home)
			if err != nil {
				Fail(licensorHint)
			}
			var term interface{}
			if *forever {
				term = "forever"
			} else {
				term = *days
			}
			bytes, err := api.Waive(licensor, *projectID, *beneficiary, *jurisdiction, term)
			if err != nil {
				Fail(err.Error())
			}
			os.Stdout.Write(bytes)
			os.Exit(0)
		}
	},
}

func waiveUsage() {
	usage := waiveDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero waive --project ID --beneficiary NAME --jurisdiction CODE (--days DAYS | --forever)\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"project ID":        projectIDLine,
			"beneficiary NAME":  "Beneficiary legal name.",
			"days DAYS":         "Term, in days.",
			"forever":           "Infinite term.",
			"jurisdiction CODE": "Beneficiary jurisdiction (ISO 3166-2, like \"US-CA\").",
		})
	Fail(usage)
}
