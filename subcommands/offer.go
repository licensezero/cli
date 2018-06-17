package subcommands

import "flag"
import "fmt"
import "os"

var Offer = Subcommand{
	Description: "Offer private licenses for sale.",
	Handler: func(args []string) {
		flagSet := flag.NewFlagSet("offer", flag.ExitOnError)
		relicense := flagSet.Int("relicense", 0, "Relicense price, in cents.")
		flagSet.Parse(args)
		if flagSet.NArg() != 1 {
			offerUsage()
		} else {
			price := args[0]
			fmt.Println(price)
			if *relicense > 0 {
				fmt.Println(*relicense)
			}
			os.Exit(0)
		}
	},
}

func offerUsage() {
	os.Stderr.WriteString(`Usage:
	<price> [--relicense CENTS]

Options:
	--relicense		Relicense price, in cents.
`)
}
