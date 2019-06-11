package subcommands

import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "licensezero.com/cli/inventory"
import "io/ioutil"
import "os"

const buyDescription = "Buy missing private licenses."

// Buy opens a buy page on licensezero.com.
var Buy = &Subcommand{
	Tag:         "buyer",
	Description: buyDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("buy", flag.ExitOnError)
		doNotOpen := doNotOpenFlag(flagSet)
		noNoncommercial := noNoncommercialFlag(flagSet)
		noProsperity := noProsperityFlag(flagSet)
		noncommercial := noncommercialFlag(flagSet)
		noReciprocal := noReciprocalFlag(flagSet)
		open := openFlag(flagSet)
		noParity := noParityFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = buyUsage
		flagSet.Parse(args)
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			Fail(identityHint)
		}
		suppressNoncommercial := *noncommercial || *noNoncommercial || *noProsperity
		suppressReciprocal := *open || *noReciprocal || *noParity
		projects, err := inventory.Inventory(paths.Home, paths.CWD, suppressNoncommercial, suppressReciprocal)
		if err != nil {
			Fail("Could not read dependency tree.")
		} else {
			licensable := projects.Licensable
			unlicensed := projects.Unlicensed
			if len(licensable) == 0 {
				os.Stdout.WriteString("No License Zero dependencies found.\n")
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
				Fail("Error sending buy request: " + err.Error())
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
				"noncommercial": noncommercialLine,
				"open":          openLine,
				"do-not-open":   doNotOpenLine,
				"silent":        silentLine,
			})
	Fail(usage)
}
