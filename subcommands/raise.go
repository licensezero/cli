package subcommands

import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const raiseDescription = "Raise Artless Devices' commission."
const commissionLine = "Agent's commission (percent)."

// Raise changes pricing.
var Raise = &Subcommand{
	Description: raiseDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("raise", flag.ExitOnError)
		commission := flagSet.Uint("commission", 0, commissionLine)
		id := idFlag(flagSet)
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = raiseUsage
		flagSet.Parse(args)
		if *commission == 0 || *id == "" {
			raiseUsage()
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
		err = api.Raise(developer, *id, *commission)
		if err != nil {
			Fail("Error sending raise request:" + err.Error())
		}
		if !*silent {
			os.Stdout.WriteString("Done.\n")
		}
		os.Exit(0)
	},
}

func raiseUsage() {
	usage := raiseDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero raise --id ID --commission PERCENT\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"commission PERCENT": commissionLine,
			"id ID":              idLine,
			"silent":             silentLine,
		})
	Fail(usage)
}
