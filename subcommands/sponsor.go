package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

const sponsorDescription = "Sponsor relicensing of a project onto permissive terms."

var Sponsor = Subcommand{
	Description: sponsorDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("sponsor", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
		projectID := ProjectID(flagSet)
		flagSet.Usage = sponsorUsage
		flagSet.Parse(args)
		if *projectID == "" {
			sponsorUsage()
		}
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			os.Stderr.WriteString(identityHint + "\n")
			os.Exit(1)
		}
		location, err := api.Sponsor(identity, *projectID)
		if err != nil {
			os.Stdout.WriteString(err.Error())
			os.Exit(1)
		}
		openURLAndExit(location, doNotOpen)
	},
}

func sponsorUsage() {
	usage := sponsorDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero sponsor --project-id ID\n\n" +
		"Options:\n" +
		"  --project-id ID  Project ID (UUID)."
	os.Stdout.WriteString(usage)
	os.Exit(1)
}
