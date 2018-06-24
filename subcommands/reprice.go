package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement reprice subcommand.

var Reprice = Subcommand{
	Description: "Change project pricing.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("reprice", flag.ContinueOnError)
		relicense := flagSet.Uint("relicense", 0, "Relicense price, in cents.")
		err := flagSet.Parse(args)
		if err != nil {
			repriceUsage()
		}
		if flagSet.NArg() != 1 {
			repriceUsage()
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

func repriceUsage() {
	os.Stderr.WriteString(`Change project pricing.

Usage:
	<price> [--relicense CENTS]

Options:
	--relicense		Relicense price, in cents.
`)
}
