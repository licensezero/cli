package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"
import "strconv"

const waiveDescription = "Generate a signed waiver."

var Waive = Subcommand{
	Description: waiveDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("waive", flag.ExitOnError)
		jurisdiction := flagSet.String("jurisdiction", "", "Jurisdiction.")
		days := flagSet.Uint("days", 0, "Days.")
		forever := flagSet.Bool("forever", false, "Forever.")
		beneficiary := flagSet.String("beneficiary", "", "Beneficiary legal name.")
		projectID := ProjectID(flagSet)
		flagSet.Usage = waiveUsage
		flagSet.Parse(args)
		if *projectID == "" {
			waiveUsage()
		} else if *forever && *days > 0 {
			waiveUsage()
		} else if *days == 0 && !*forever {
			waiveUsage()
		} else if *beneficiary == "" || *jurisdiction == "" {
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
	usage := waiveDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero waive --project-id ID --beneficiary NAME --jurisdiction CODE (--days DAYS | --forever)\n\n" +
		"Options:\n" +
		"  --project-id ID      " + projectIDLine + "\n" +
		"  --beneficiary NAME   Beneficiary legal name.\n" +
		"  --days DAYS          Term, in days.\n" +
		"  --forever            Infinite term.\n" +
		"  --jurisdiction CODE  Beneficiary jurisdiction (ISO 3166-2, like \"US-CA\").\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
