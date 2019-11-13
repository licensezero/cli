package subcommands

import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const waiveDescription = "Generate a waiver."

// Waive generates a signed waiver.
var Waive = &Subcommand{
	Tag:         "seller",
	Description: waiveDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("waive", flag.ExitOnError)
		jurisdiction := flagSet.String("jurisdiction", "", "Jurisdiction.")
		days := flagSet.Uint("days", 0, "Days.")
		forever := flagSet.Bool("forever", false, "Forever.")
		beneficiary := flagSet.String("beneficiary", "", "Beneficiary legal name.")
		offerID := offerIDFlag(flagSet)
		id := idFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = waiveUsage
		flagSet.Parse(args)
		if *offerID == "" && *id == "" {
			licenseUsage()
		} else if *offerID != "" && *id != "" {
			licenseUsage()
		} else if *forever && *days > 0 {
			waiveUsage()
		} else if *days == 0 && !*forever {
			waiveUsage()
		} else if *beneficiary == "" || *jurisdiction == "" {
			waiveUsage()
		}
		if *offerID != "" {
			*id = *offerID
		}
		if !validID(*id) {
			invalidID()
		}
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
		bytes, err := api.Waive(licensor, *id, *beneficiary, *jurisdiction, term)
		if err != nil {
			Fail("Error sending waiver request: " + err.Error())
		}
		os.Stdout.Write(bytes)
		os.Exit(0)
	},
}

func waiveUsage() {
	usage := waiveDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero waive --id ID --beneficiary NAME --jurisdiction CODE (--days DAYS | --forever)\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"id ID":             idLine,
			"beneficiary NAME":  "Beneficiary legal name.",
			"days DAYS":         "Term, in days.",
			"forever":           "Infinite term.",
			"jurisdiction CODE": "Beneficiary jurisdiction (ISO 3166-2, like \"US-CA\").",
		})
	Fail(usage)
}
