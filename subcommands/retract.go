package subcommands

import "os"
import "fmt"

// TODO: Implement retract subcommand.

var Retract = Subcommand{
	Description: "Retract a package from sale.",
	Handler: func(args []string, paths Paths) {
		if len(args) != 1 {
			os.Stderr.WriteString("<project id>")
			os.Exit(1)
		} else {
			projectID := args[0]
			fmt.Println(projectID)
			os.Exit(0)
		}
	},
}
