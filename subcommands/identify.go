package subcommands

import "os"
import "fmt"

func Identify(args []string) {
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
}
