package subcommands

import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const freebieDescription = "Generate a waiver."

// Freebie generates a signed waiver.
var Freebie = &Subcommand{
	Description: freebieDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("freebie", flag.ExitOnError)
		days := flagSet.Uint("days", 0, "Days.")
		forever := flagSet.Bool("forever", false, "Forever.")
		name := flagSet.String("name", "", "User Legal Name.")
		email := flagSet.String("email", "", "User E-Mail.")
		jurisdiction := flagSet.String("jurisdiction", "", "User Jurisdiction.")
		offerID := offerIDFlag(flagSet)
		id := idFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = freebieUsage
		flagSet.Parse(args)
		if *offerID == "" && *id == "" {
			freebieUsage()
		} else if *offerID != "" && *id != "" {
			freebieUsage()
		} else if *forever && *days > 0 {
			freebieUsage()
		} else if *days == 0 && !*forever {
			freebieUsage()
		} else if *name == "" || *jurisdiction == "" || *email == "" {
			freebieUsage()
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
		var term interface{}
		if *forever {
			term = "forever"
		} else {
			term = *days
		}
		bytes, err := api.Freebie(developer, *id, *name, *jurisdiction, *email, term)
		if err != nil {
			Fail("Error sending waiver request: " + err.Error())
		}
		os.Stdout.Write(bytes)
		os.Exit(0)
	},
}

func freebieUsage() {
	usage := freebieDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero freebie --id ID --name NAME --email EMAIL --jurisdiction CODE (--days DAYS | --forever)\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"id ID":             idLine,
			"name NAME":         "User legal name.",
			"email EMAIL":       "User e-mail.",
			"jurisdiction CODE": "User jurisdiction (ISO 3166-2, like \"US-CA\").",
			"days DAYS":         "Term, in days.",
			"forever":           "Infinite term.",
		})
	Fail(usage)
}
