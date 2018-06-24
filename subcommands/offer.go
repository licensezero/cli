package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement offer subcommand.

const offerDescription = "Offer private licenses for sale."

var Offer = Subcommand{
	Description: offerDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("offer", flag.ContinueOnError)
		relicense := Relicense(flagSet)
		price := Price(flagSet)
		err := flagSet.Parse(args)
		if err != nil {
			offerUsage()
		}
		if *price == 0 {
			offerUsage()
		}
		fmt.Println(price)
		if *relicense > 0 {
			fmt.Println(*relicense)
		}
		os.Exit(0)
	},
}

func offerUsage() {
	usage := offerDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero offer --price CENTS [--relicense CENTS]\n\n" +
		"Options:\n" +
		"  --price CENTS      " + priceLine + "\n" +
		"  --relicense CENTS  " + relicenseLine + "\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
