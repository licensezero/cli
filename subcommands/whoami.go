package subcommands

import "fmt"
import "github.com/licensezero/cli/data"
import "os"

const whoAmIDescription = "Show your licensee identity."

var WhoAmI = Subcommand{
	Description: whoAmIDescription,
	Handler: func(args []string, paths Paths) {
		data, err := data.ReadIdentity(paths.Home)
		if err != nil {
			os.Stderr.WriteString("Could not read identity file.\n")
			os.Exit(1)
		} else {
			fmt.Println("Name: " + data.Name)
			fmt.Println("Jurisdiction: " + data.Jurisdiction)
			fmt.Println("E-Mail: " + data.EMail)
			os.Exit(0)
		}
	},
}
