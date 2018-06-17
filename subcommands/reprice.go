package subcommands

import "flag"
import "fmt"
import "os"

func Reprice(args []string) {
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
}

func repriceUsage() {
	os.Stderr.WriteString(`Change project pricing.

Usage:
	<price> [--relicense CENTS]

Options:
	--relicense		Relicense price, in cents.
`)
}
