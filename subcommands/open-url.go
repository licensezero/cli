package subcommands

import "github.com/skratchdot/open-golang/open"
import "os"

func openURLAndExit(url string, noBrowser *bool) {
	os.Stdout.WriteString(url + "\n")
	if !*noBrowser {
		open.Run(url)
	}
	os.Exit(0)
}
