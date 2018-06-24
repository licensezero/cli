package subcommands

import "flag"
import "fmt"
import "os"

// TODO: Implement retract subcommand.

const retractDescription = "Retract a package from sale."

var Retract = Subcommand{
	Description: retractDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("retract", flag.ExitOnError)
		projectID := ProjectID(flagSet)
		flagSet.Usage = retractUsage
		flagSet.Parse(args)
		if *projectID == "" {
			retractUsage()
		}
		if len(args) != 1 {
			retractUsage()
		}
		fmt.Println(projectID)
		os.Exit(0)
	},
}

func retractUsage() {
	usage := retractDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero retract --project-id ID\n\n" +
		"Options:\n" +
		"  --project-id ID  " + projectIDLine + "\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
