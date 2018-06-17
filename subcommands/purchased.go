package subcommands

import "os"
import "fmt"

func Purchased(args []string) {
	if len(args) != 1 {
		os.Stderr.WriteString("<URL>")
		os.Exit(1)
	} else {
		url := args[0]
		fmt.Println(url)
		os.Exit(0)
	}
}
