package subcommands

import "flag"
import "fmt"
import "github.com/licensezero/cli/inventory"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "github.com/skratchdot/open-golang/open"
import "os"

var Buy = Subcommand{
	Description: "Buy missing private licenses.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("buy", flag.ContinueOnError)
		doNotOpen := DoNotOpen(flagSet)
		noNoncommercial := flagSet.Bool("no-noncommercial", false, "Ignore L0-NC dependencies.")
		noReciprocal := flagSet.Bool("no-reciprocal", false, "Ignore L0-R dependencies.")
		err := flagSet.Parse(args)
		if err != nil {
			os.Stderr.WriteString(`Buy missing private licenses.

Options;
	--no-noncommercial  Ignore packages under noncommerical terms.
	--no-reciprocal     Ignore packages under reciprocal terms.
	--do-not-open       Do not open buy page in web browser.
`)
			os.Exit(1)
		}
		identity, err := data.ReadIdentity(paths.Home)
		if err != nil {
			os.Stderr.WriteString("Create an identity with `licensezero identify` first.")
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
			os.Stdout.WriteString(location + "\n")
			if !*doNotOpen {
				open.Run(location)
			}
			os.Exit(0)
		}
	},
}
