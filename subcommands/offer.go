package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement offer subcommand.

var Offer = Subcommand{
	Description: "Offer private licenses for sale.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("offer", flag.ContinueOnError)
		relicense := flagSet.Int("relicense", 0, "Relicense price, in cents.")
		err := flagSet.Parse(args)
		if err != nil {
			offerUsage()
		}
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
	os.Stderr.WriteString(`Offer private licenses for sale.

Usage:
	<price> [--relicense CENTS]

Options:
	--relicense		Relicense price, in cents.
`)
}
