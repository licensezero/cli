package subcommands

import "os"
import "fmt"

var Import = Subcommand{
	Description: "Import a private license or waiver from file.",
	Handler: func(args []string, home string) {
		if len(args) != 1 {
			os.Stderr.WriteString("<file>\n")
			os.Exit(1)
		} else {
			file := args[0]
			fmt.Println(file)
			os.Exit(0)
		}
	},
}
