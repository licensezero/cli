package subcommands

import "flag"
import "fmt"
import "os"

var Reprice = Subcommand{
	Description: "Change project pricing.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("reprice", flag.ExitOnError)
		relicense := flagSet.Int("relicense", 0, "Relicense price, in cents.")
		flagSet.Parse(args)
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
	os.Stderr.WriteString(`Usage:
	<price> [--relicense CENTS]

Options:
	--relicense		Relicense price, in cents.
`)
}
