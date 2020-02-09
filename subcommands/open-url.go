package subcommands

import (
	"github.com/skratchdot/open-golang/open"
	"io"
)

func openURL(url string, noBrowser *bool, stdout io.StringWriter) {
	stdout.WriteString(url + "\n")
	if !*noBrowser {
		open.Run(url)
	}
}
