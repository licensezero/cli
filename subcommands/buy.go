package subcommands

import "flag"
import "fmt"
import "github.com/licensezero/cli/inventory"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

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
			os.Stderr.WriteString("Could not read dependeny tree.")
			os.Exit(1)
		} else {
			licensable := projects.Licensable
			unlicensed := projects.Unlicensed
			if len(licensable) == 0 {
				fmt.Println("No License Zero depedencies found.")
				os.Exit(0)
			}
			if len(unlicensed) == 0 {
				fmt.Println("No private licenses to buy.")
				os.Exit(0)
			}
			var projectIDs []string
			for _, project := range unlicensed {
				projectIDs = append(projectIDs, project.Envelope.Manifest.ProjectID)
			}
			location, err := api.Buy(identity, projectIDs)
			if err != nil {
				fmt.Println(err.Error())
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
			"  --no-noncommercial  " + noNoncommercialLine + "\n" +
			"  --no-reciprocal     " + noReciprocalLine + "\n" +
			"  --do-not-open       " + doNotOpenLine + "\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
