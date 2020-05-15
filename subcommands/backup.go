package subcommands

import "time"
import "licensezero.com/cli/data"
import "github.com/mholt/archiver"
import "os"

const backupDescription = "Create a tarball of your data."

// Backup writes a tarball of configuration files to the current directory.
var Backup = &Subcommand{
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
