package subcommands

import "flag"
import "fmt"
import "os"

var Buy = Subcommand{
	Description: "Buy missing private licenses.",
	Handler: func(args []string, home string) {
		flagSet := flag.NewFlagSet("buy", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
		flagSet.Parse(args)
		if *doNotOpen {
			fmt.Println("not opening")
		}
		os.Exit(0)
	},
}
