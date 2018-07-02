package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "io/ioutil"

const sponsorDescription = "Sponsor relicensing of a project."

// Sponsor starts a project sponsorship transaction.
var Sponsor = Subcommand{
	Tag:         "buyer",
	Description: sponsorDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("sponsor", flag.ExitOnError)
		doNotOpen := doNotOpenFlag(flagSet)
		projectID := projectIDFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = sponsorUsage
		flagSet.Parse(args)
		if *projectID == "" {
			sponsorUsage()
		}
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			Fail(identityHint)
		}
		location, err := api.Sponsor(identity, *projectID)
		if err != nil {
			Fail("Error sending sponsor request: " + err.Error())
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
	Fail(usage)
}
