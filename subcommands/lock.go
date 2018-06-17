package subcommands

import "os"
import "fmt"

func Lock(args []string) {
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
}
