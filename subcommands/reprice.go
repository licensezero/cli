package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement reprice subcommand.

const repriceDescription = "Change project pricing"

var Reprice = Subcommand{
	Description: repriceDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("reprice", flag.ExitOnError)
		price := Price(flagSet)
		relicense := Relicense(flagSet)
		flagSet.Usage = repriceUsage
		flagSet.Parse(args)
		if flagSet.NArg() != 1 {
			repriceUsage()
		} else {
			fmt.Println(price)
			if *relicense > 0 {
				fmt.Println(*relicense)
			}
			os.Exit(0)
		}
	},
}

func repriceUsage() {
	usage := repriceDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero reprice --price CENTS [--relicense CENTS]\n\n" +
		"Options:\n" +
		"  --price      " + priceLine + "\n" +
		"  --relicense  " + relicenseLine + "\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
