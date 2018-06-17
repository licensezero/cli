package subcommands

import "os"
import "fmt"

func Import(args []string) {
	if len(args) != 1 {
		os.Stderr.WriteString("<file>\n")
		os.Exit(1)
	} else {
		file := args[0]
		fmt.Println(file)
		os.Exit(0)
	}
}
