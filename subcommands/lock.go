package subcommands

import "os"
import "fmt"

var Lock = Subcommand{
	Description: "Lock project pricing",
	Handler: func(args []string, home string) {
		if len(args) != 2 {
			os.Stderr.WriteString("<project id> <date>\n")
			os.Exit(1)
		} else {
			projectID := args[0]
			date := args[1]
			fmt.Println(projectID)
			fmt.Println(date)
			os.Exit(0)
		}
	},
}
