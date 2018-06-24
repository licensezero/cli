package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

const lockDescription = "Lock project pricing."

var Lock = Subcommand{
	Description: lockDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("lock", flag.ExitOnError)
		projectID := ProjectID(flagSet)
		unlock := flagSet.String("unlock", "", "")
		flagSet.Usage = lockUsage
		flagSet.Parse(args)
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			os.Stderr.WriteString(licensorHint + "\n")
			os.Exit(1)
		}
		err = api.Lock(licensor, *projectID, *unlock)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func lockUsage() {
	usage := lockDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero lock --project-id ID --unlock DATE\n\n" +
		"Options:\n" +
		"  --project-id  " + projectIDLine + "\n" +
		"  --unlock      Unlock date.\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
