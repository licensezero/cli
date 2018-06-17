package subcommands

import "os"
import "fmt"

var Purchased = Subcommand{
	Description: "Import a bundle of purchased licenses from URL.",
	Handler: func(args []string, paths Paths) {
		if len(args) != 1 {
			os.Stderr.WriteString("<URL>")
			os.Exit(1)
		} else {
			url := args[0]
			fmt.Println(url)
			os.Exit(0)
		}
	},
}
