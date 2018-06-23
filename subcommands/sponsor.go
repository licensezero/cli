package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement sponsor subcommand.

var Sponsor = Subcommand{
	Description: "Sponsor relicensing of a project onto permissive terms.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("sponsor", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
		flagSet.Parse(args)
		if *doNotOpen {
			fmt.Println("not opening")
		}
		os.Exit(0)
	},
}
