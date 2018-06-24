package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement purchased subcommand.

const purchasedDescription = "Import a bundle of purchased licenses from URL."

var Purchased = Subcommand{
	Description: purchasedDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("purchased", flag.ContinueOnError)
		bundle := flagSet.String("bundle", "", "")
		err := flagSet.Parse(args)
		if err != nil || *bundle == "" {
			purchasedUsage()
		}
		fmt.Println(bundle)
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
