package subcommands

import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

const lockDescription = "Lock project pricing and availability until a given date."

var Lock = Subcommand{
	Description: lockDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("lock", flag.ExitOnError)
		projectID := ProjectID(flagSet)
		unlock := flagSet.String("unlock", "", "")
		silent := Silent(flagSet)
		flagSet.Usage = lockUsage
		flagSet.Parse(args)
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			os.Stderr.WriteString(licensorHint + "\n")
			os.Exit(1)
		}
		err = api.Lock(licensor, *projectID, *unlock)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
		if !*silent {
			os.Stdout.WriteString("Locked pricing.\n")
		}
		os.Exit(0)
	},
}

func lockUsage() {
	usage := lockDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero lock --project ID --unlock DATE\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"project": projectIDLine,
			"silent":  silentLine,
			"unlock":  "Unlock date.",
		})
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
