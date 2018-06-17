package subcommands

import "flag"
import "fmt"
import "os"

var Waive = Subcommand{
	Description: "Generate a signed waiver.",
	Handler: func(args []string, home string) {
		flagSet := flag.NewFlagSet("waive", flag.ExitOnError)
		jurisdiction := flagSet.Bool("jurisdiction", false, "Jurisdiction.")
		days := flagSet.Int("days", 0, "Days.")
		forever := flagSet.Bool("forever", false, "Forever.")
		beneficiary := flagSet.Bool("beneficiary", false, "Beneficiary legal name.")
		flagSet.Parse(args)
		if *forever && *days > 0 {
			waiveUsage()
		} else if flagSet.NArg() != 1 {
			waiveUsage()
		} else {
			projectID := args[0]
			fmt.Println(projectID)
			fmt.Println(beneficiary)
			fmt.Println(jurisdiction)
			fmt.Println(days)
			os.Exit(0)
		}
	},
}

func waiveUsage() {
	os.Stderr.WriteString(`Usage:
	<project id> --beneficiary NAME --jurisdiction CODE (--days DAYS | --forever)

Options:
	--days DAYS          Term in days.
	--forever            Infinite term.
	--jurisdiction CODE  Beneficiary jurisdiction (ISO 3166-2).
	--beneficiary NAME   Beneficiary name.
`)
	os.Exit(1)
}
