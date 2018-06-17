package subcommands

import "os"
import "fmt"

var WhoAmI = Subcommand{
	Description: "Show your licensee identity.",
	Handler: func(args []string, home string) {
		data, err := readIdentity(home)
		if err != nil {
			os.Stderr.WriteString("Could not read identity file.\n")
			os.Exit(1)
		} else {
			fmt.Println("Name: " + data.Name)
			fmt.Println("Jurisdiction: " + data.Jurisdiction)
			fmt.Println("E-Mail: " + data.Email)
			os.Exit(0)
		}
	},
}
