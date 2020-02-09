package subcommands

import (
	"io/ioutil"
	"net/http"
	"os"
)

const latestDescription = "Check for a newer version."

// Latest prints checks the running version against the latest available.
var Latest = &Subcommand{
	Tag:         "misc",
	Description: latestDescription,
	Handler: func(args []string, stdin, stdout, stderr *os.File) int {
		var running string
		if args[0] == "" {
			running = "Development Build"
		} else {
			running = "v" + args[0]
		}
		response, err := http.Get("https://licensezero.com/cli-version")
		if err != nil {
			stderr.WriteString("Could not fetch latest version from licensezero.com.\n")
			return 1
		}
		defer response.Body.Close()
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			stderr.WriteString("Error reading response body.\n")
			return 1
		}
		current := string(responseBody)
		stdout.WriteString("Running: " + running + "\n")
		stdout.WriteString("Latest:  " + current + "\n")
		if running == current {
			return 0
		}
		response, err = http.Get("https://licensezero.com/one-line-install.sh")
		if err != nil {
			return 1
		}
		defer response.Body.Close()
		responseBody, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return 1
		}
		stdout.WriteString("Install: " + string(responseBody) + "\n")
		return 1
	},
}
