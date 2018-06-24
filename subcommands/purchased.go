package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement purchased subcommand.

const purchasedDescription = "Download a bundle of purchased licenses."

var Purchased = Subcommand{
	Description: purchasedDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("purchased", flag.ExitOnError)
		bundle := flagSet.String("bundle", "", "")
		flagSet.Usage = purchasedUsage
		flagSet.Parse(args)
		if *bundle == "" {
			purchasedUsage()
		}
		fmt.Println(*bundle)
		os.Exit(0)
	},
}

func purchasedUsage() {
	usage := purchasedDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero purchased --bundle URL\n\n" +
		"Options:\n" +
		"  --bundle URL  URL of purchase bundle to import.\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
