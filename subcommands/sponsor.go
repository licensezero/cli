package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement sponsor subcommand.

const sponsorDescription = "Sponsor relicensing of a project onto permissive terms."

var Sponsor = Subcommand{
	Description: sponsorDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("sponsor", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
		flagSet.Usage = sponsorUsage
		flagSet.Parse(args)
		if *doNotOpen {
			fmt.Println("not opening")
		}
		os.Exit(0)
	},
}

func sponsorUsage() {
	usage := sponsorDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero sponsor\n"
	os.Stdout.WriteString(usage)
	os.Exit(1)
}
