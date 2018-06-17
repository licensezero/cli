package subcommands

import "flag"
import "fmt"
import "os"

func Quote(args []string) {
	flagSet := flag.NewFlagSet("quote", flag.ExitOnError)
	noNoncommercial := flagSet.Bool("no-noncommercial", false, "Ignore L0-NC dependencies.")
	noReciprocal := flagSet.Bool("no-reciprocal", false, "Ignore L0-R dependencies.")
	flagSet.Parse(args)
	if *noNoncommercial {
		fmt.Println("No L0-NC")
	}
	if *noReciprocal {
		fmt.Println("No L0-R")
	}
	os.Exit(0)
}