package subcommands

import "os"
import "fmt"

var Identify = Subcommand{
	Description: "Identify yourself.",
	Handler: func(args []string, home string) {
		if len(args) != 3 {
			os.Stderr.WriteString("<name> <jurisdiction> <email>\n")
			os.Exit(1)
		} else {
			name := args[0]
			jurisdiction := args[1]
			email := args[2]
			fmt.Println(name)
			fmt.Println(jurisdiction)
			fmt.Println(email)
			os.Exit(0)
		}
	},
}
