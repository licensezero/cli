package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"
import "strconv"

var Waive = Subcommand{
	Description: "Generate a signed waiver.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("waive", flag.ContinueOnError)
		jurisdiction := flagSet.String("jurisdiction", "", "Jurisdiction.")
		days := flagSet.Uint("days", 0, "Days.")
		forever := flagSet.Bool("forever", false, "Forever.")
		beneficiary := flagSet.String("beneficiary", "", "Beneficiary legal name.")
		err := flagSet.Parse(args)
		if err != nil {
			waiveUsage()
		}
		if *forever && *days > 0 {
			waiveUsage()
		} else if *days == 0 && !*forever {
			waiveUsage()
		} else if flagSet.NArg() != 1 || *beneficiary == "" || *jurisdiction == "" {
			waiveUsage()
		} else {
			projectID := args[0]
			licensor, err := data.ReadLicensor(paths.Home)
			if err != nil {
				os.Stderr.WriteString("Create a licensor identity with `licensezero register` or `licensezero set-licensor-id`.")
				os.Exit(1)
			}
			var term string
			if *forever {
				term = "forever"
			} else {
				term = strconv.Itoa(int(*days))
			}
			bytes, err := api.Waive(licensor, projectID, *beneficiary, *jurisdiction, term)
			if err != nil {
				os.Stderr.WriteString(err.Error())
				os.Exit(1)
			}
			os.Stdout.Write(bytes)
			os.Exit(0)
		}
	},
}

func waiveUsage() {
	os.Stderr.WriteString(`Generate a signed waiver.
Usage:
	<project id> --beneficiary NAME --jurisdiction CODE (--days DAYS | --forever)

Options:
	--days DAYS          Term in days.
	--forever            Infinite term.
	--jurisdiction CODE  Beneficiary jurisdiction (ISO 3166-2).
	--beneficiary NAME   Beneficiary name.
`)
	os.Exit(1)
}
