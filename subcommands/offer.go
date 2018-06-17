package subcommands

import "flag"
import "fmt"
import "os"

func Offer(args []string) {
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
}

func offerUsage() {
	os.Stderr.WriteString(`Offer private licenses for sale.

Usage:
	<price> [--relicense CENTS]

Options:
	--relicense		Relicense price, in cents.
`)
}
