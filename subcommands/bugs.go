package subcommands

import (
	"flag"
	"io"
	"io/ioutil"
	"net/http"
)

const bugsDescription = "Open the CLI bug tracker page."

// Bugs opens the CLI tracker bug tracker page.
var Bugs = &Subcommand{
	Tag:         "misc",
	Description: bugsDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		flagSet := flag.NewFlagSet("bugs", flag.ExitOnError)
		doNotOpen := doNotOpenFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			usage := bugsDescription + "\n\n" +
				"Usage:\n" +
				"  licensezero bugs\n\n" +
				"Options:\n" +
				flagsList(map[string]string{
					"do-not-open": doNotOpenUsage,
				})
			stderr.WriteString(usage)
		}
		err := flagSet.Parse(args)
		if err != nil {
			return 1
		}
		openURL("https://github.com/licensezero/cli/issues", doNotOpen, stdout)
		return 0
	},
}
