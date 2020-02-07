package subcommands

import (
	"github.com/mholt/archiver/v3"
	"licensezero.com/licensezero/user"
	"os"
	"time"
)

const backupDescription = "Create a tarball of your data."

// Backup writes a tarball of configuration files to the current directory.
var Backup = &Subcommand{
	Tag:         "misc",
	Description: backupDescription,
	Handler: func(args []string, paths Paths) {
		now := time.Now()
		fileName := "licensezero-backup-" + now.Format(time.RFC3339) + ".tar"
		tar := archiver.NewTar()
		err := tar.Archive([]string{user.ConfigPath(paths.Home)}, fileName)
		if err != nil {
			Fail("Error creating tarball.")
		}
		os.Exit(0)
	},
}
