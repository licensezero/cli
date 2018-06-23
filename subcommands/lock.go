package subcommands

import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "os"

var Lock = Subcommand{
	Description: "Lock project pricing",
	Handler: func(args []string, paths Paths) {
		if len(args) != 2 {
			os.Stderr.WriteString(`Lock project pricing.

Usage;
  <project id> <date>
`)
			os.Exit(1)
		} else {
			projectID := args[0]
			date := args[1]
			licensor, err := data.ReadLicensor(paths.Home)
			if err != nil {
				os.Stderr.WriteString("Create a licensor identity with `licensezero register` or `licensezero set-licensor-id`.")
				os.Exit(1)
			}
			err = api.Lock(licensor, projectID, date)
			if err != nil {
				os.Stderr.WriteString(err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		}
	},
}
