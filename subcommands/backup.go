package subcommands

import (
	"github.com/mholt/archiver/v3"
	"licensezero.com/licensezero/user"
	"time"
)

const backupDescription = "Create a tarball of your data."

// Backup writes a tarball of configuration files to the current directory.
var Backup = &Subcommand{
	Tag:         "misc",
	Description: backupDescription,
	Handler: func(env Environment) int {
		now := time.Now()
		fileName := "licensezero-backup-" + now.Format(time.RFC3339) + ".tar"
		tar := archiver.NewTar()
		configPath, err := user.ConfigPath()
		if err != nil {
			env.Stderr.WriteString("Could not compute configuration path.")
			return 1
		}
		err = tar.Archive([]string{configPath}, fileName)
		if err != nil {
			env.Stderr.WriteString("Error creating tarball.")
			return 1
		}
		return 0
	},
}
