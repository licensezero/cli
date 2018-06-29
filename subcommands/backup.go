package subcommands

import "time"
import "github.com/licensezero/cli/data"
import "github.com/mholt/archiver"
import "os"

const backupDescription = "Create a tarball of your data."

var Backup = Subcommand{
	Tag:         "misc",
	Description: backupDescription,
	Handler: func(args []string, paths Paths) {
		now := time.Now()
		fileName := "licensezero-backup-" + now.Format(time.RFC3339) + ".tar"
		err := archiver.Tar.Make(fileName, []string{data.ConfigPath(paths.Home)})
		if err != nil {
			Fail("Error creating tarball.")
		}
		os.Exit(0)
	},
}
