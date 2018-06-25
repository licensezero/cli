package subcommands

import "flag"
import "github.com/licensezero/cli/inventory"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

// TODO: licensezero buy --json

const buyDescription = "Buy private licenses you are missing."

var Buy = Subcommand{
	Description: buyDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("buy", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
		noNoncommercial := NoNoncommercial(flagSet)
		noReciprocal := NoReciprocal(flagSet)
		flagSet.Usage = buyUsage
		flagSet.Parse(args)
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			os.Stderr.WriteString(identityHint + "\n")
			os.Exit(1)
		}
		projects, err := inventory.Inventory(paths.Home, paths.CWD, *noNoncommercial, *noReciprocal)
		if err != nil {
			os.Stderr.WriteString("Could not read dependeny tree.\n")
			os.Exit(1)
		} else {
			licensable := projects.Licensable
			unlicensed := projects.Unlicensed
			if len(licensable) == 0 {
				os.Stdout.WriteString("No License Zero depedencies found.\n")
				os.Exit(0)
			}
			if len(unlicensed) == 0 {
				os.Stdout.WriteString("No private licenses to buy.\n")
				os.Exit(0)
			}
			var projectIDs []string
			for _, project := range unlicensed {
				projectIDs = append(projectIDs, project.Envelope.Manifest.ProjectID)
			}
			location, err := api.Buy(identity, projectIDs)
			if err != nil {
				os.Stderr.WriteString(err.Error() + "\n")
				os.Exit(1)
			}
			openURLAndExit(location, doNotOpen)
		}
	},
}

func buyUsage() {
	usage :=
		buyDescription + "\n\n" +
			"Usage:\n" +
			"  licensezero buy\n\n" +
			"Options:\n" +
			flagsList(map[string]string{
				"no-noncommercial": noNoncommercialLine,
				"no-reciprocal":    noReciprocalLine,
				"do-not-open":      doNotOpenLine,
				"silent":           silentLine,
			})
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
