package subcommands

import "flag"
import "fmt"
import "os"

func Buy(args []string) {
	flagSet := flag.NewFlagSet("buy", flag.ExitOnError)
	doNotOpen := flagSet.Bool("do-not-open", false, "Do not open checkout page.")
	flagSet.Parse(args)
	if *doNotOpen {
		fmt.Println("not opening")
	}
	os.Exit(0)
}
