package subcommands

import "os"
import "fmt"

var SetLicensorID = Subcommand{
	Description: "Set your licensezero.com licensor ID.",
	Handler: func(args []string) {
		if len(args) != 1 {
			os.Stderr.WriteString("Usage:\n\t<licensor ID>\n")
			os.Exit(1)
		} else {
			licensorID := args[0]
			fmt.Println(licensorID)
			os.Exit(0)
		}
	},
}
