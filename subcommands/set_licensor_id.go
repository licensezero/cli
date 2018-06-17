package subcommands

import "os"
import "fmt"

func SetLicensorID(args []string) {
	if len(args) != 1 {
		os.Stderr.WriteString(`Set your License Zero licensor ID.

Usage:
	<licensor ID>
`)
		os.Exit(1)
	} else {
		licensorID := args[0]
		fmt.Println(licensorID)
		os.Exit(0)
	}
}
