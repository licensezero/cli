package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement sponsor subcommand.

var Sponsor = Subcommand{
	Description: "Sponsor relicensing of a project onto permissive terms.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("sponsor", flag.ContinueOnError)
		doNotOpen := DoNotOpen(flagSet)
		err := flagSet.Parse(args)
		if err != nil {
			sponsorUsage()
		}
		if *doNotOpen {
			fmt.Println("not opening")
		}
		os.Exit(0)
	},
}

func sponsorUsage() {
	os.Stdout.WriteString(`Sponsor relicensing of a project onto permissive terms.`)
	os.Exit(1)
}
