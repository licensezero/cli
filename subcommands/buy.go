package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "github.com/licensezero/cli/inventory"
import "io/ioutil"
import "os"

const buyDescription = "Buy missing private licenses."

var Buy = Subcommand{
	Tag:         "buyer",
	Description: buyDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("buy", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
		noNoncommercial := NoNoncommercial(flagSet)
		noReciprocal := NoReciprocal(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = buyUsage
		flagSet.Parse(args)
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			Fail(identityHint)
		}
		projects, err := inventory.Inventory(paths.Home, paths.CWD, *noNoncommercial, *noReciprocal)
		if err != nil {
			Fail("Could not read dependeny tree.")
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
				Fail(err.Error())
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
	Fail(usage)
}
