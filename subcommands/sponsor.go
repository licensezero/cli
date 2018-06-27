package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

const sponsorDescription = "Sponsor relicensing of a project onto Charity terms."

var Sponsor = Subcommand{
	Tag:         "buyer",
	Description: sponsorDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("sponsor", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
		projectID := ProjectID(flagSet)
		flagSet.SetOutput(ioutil.Discard)
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
			os.Stdout.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
		openURLAndExit(location, doNotOpen)
	},
}

func sponsorUsage() {
	usage := sponsorDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero sponsor --project ID\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"project ID": "Project ID (UUID).",
		})
	os.Stdout.WriteString(usage)
	os.Exit(1)
}
