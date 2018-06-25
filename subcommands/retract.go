package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

const retractDescription = "Retract project licenses for a project from sale."

var Retract = Subcommand{
	Description: retractDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("retract", flag.ExitOnError)
		projectID := ProjectID(flagSet)
		silent := Silent(flagSet)
		flagSet.Usage = retractUsage
		flagSet.Parse(args)
		if *projectID == "" {
			retractUsage()
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			os.Stderr.WriteString(licensorHint + "\n")
			os.Exit(1)
		}
		err = api.Retract(licensor, *projectID)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
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
		"  licensezero retract --project-id ID\n\n" +
		"Options:\n" +
		"  --project-id ID  " + projectIDLine + "\n" +
		"  --silent         " + silentLine + "\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
