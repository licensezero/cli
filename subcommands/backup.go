package subcommands

import "time"
import "github.com/licensezero/cli/data"
import "github.com/mholt/archiver"
import "os"

const backupDescription = "Create a dated tarball of your License Zero data."

var Backup = Subcommand{
	Tag:         "misc",
	Description: backupDescription,
	Handler: func(args []string, paths Paths) {
		now := time.Now()
		fileName := "licensezero-backup-" + now.Format(time.RFC3339) + ".tar"
		err := archiver.Tar.Make(fileName, []string{data.ConfigPath(paths.Home)})
		if err != nil {
			os.Stdout.WriteString("Error creating tarball.\n")
			os.Exit(1)
		}
		os.Exit(0)
	},
}
