package subcommands

import (
	"github.com/skratchdot/open-golang/open"
	"os"
)

func openURL(url string, noBrowser *bool, stdout *os.File) {
	stdout.WriteString(url + "\n")
	if !*noBrowser {
		open.Run(url)
	}
}
