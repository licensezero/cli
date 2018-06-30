package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

const retractDescription = "Stop offering private licenses for sale."

// Retract pulls a project from sale.
var Retract = Subcommand{
	Tag:         "seller",
	Description: retractDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("retract", flag.ExitOnError)
		projectID := projectIDFlag(flagSet)
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = retractUsage
		flagSet.Parse(args)
		if *projectID == "" {
			retractUsage()
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		err = api.Retract(licensor, *projectID)
		if err != nil {
			Fail(err.Error())
		}
		if !*silent {
			os.Stdout.WriteString("Retracted from sale.\n")
		}
		os.Exit(0)
	},
}

func retractUsage() {
	usage := retractDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero retract --project ID\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"project ID": projectIDLine,
			"silent":     silentLine,
		})
	Fail(usage)
}
