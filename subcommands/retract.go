package subcommands

import "os"
import "fmt"

func Retract(args []string) {
	if len(args) != 1 {
		os.Stderr.WriteString("<project id>")
		os.Exit(1)
	} else {
		projectID := args[0]
		fmt.Println(projectID)
		os.Exit(0)
	}
}
