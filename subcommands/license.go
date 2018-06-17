package subcommands

import "flag"
import "fmt"
import "os"

func License(args []string) {
	flagSet := flag.NewFlagSet("license", flag.ExitOnError)
	noncommercial := flagSet.Bool("noncommercial", false, "Use noncommercial public license.")
	reciprocal := flagSet.Bool("reciprocal", false, "Use reciprocal public license.")
	flagSet.Parse(args)
	if flagSet.NArg() != 1 {
		licenseUsage()
	} else {
		projectID := flagSet.Args()[0]
		if *noncommercial && *reciprocal {
			licenseUsage()
		}
		fmt.Println(projectID)
		os.Exit(0)
	}
}

func licenseUsage() {
	os.Stderr.WriteString(`Write a public license file.

Usage:
	 <project id> (--noncommercial | --reciprocal)
`)
	os.Exit(1)
}
