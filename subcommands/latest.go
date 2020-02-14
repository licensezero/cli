package subcommands

import (
	"io/ioutil"
	"net/http"
)

const latestDescription = "Check for a newer version."

// Latest prints checks the running version against the latest available.
var Latest = &Subcommand{
	Tag:         "misc",
	Description: latestDescription,
	Handler: func(env Environment) int {
		var running string
		if env.Rev == "" {
			running = "Development Build"
		} else {
			running = "v" + env.Rev
		}
		response, err := http.Get("https://licensezero.com/cli-version")
		if err != nil {
			env.Stderr.WriteString("Could not fetch latest version from licensezero.com.\n")
			return 1
		}
		defer response.Body.Close()
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			env.Stderr.WriteString("Error reading response body.\n")
			return 1
		}
		current := string(responseBody)
		env.Stdout.WriteString("Running: " + running + "\n")
		env.Stdout.WriteString("Latest:  " + current + "\n")
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
		env.Stdout.WriteString("Install: " + string(responseBody) + "\n")
		return 1
	},
}
