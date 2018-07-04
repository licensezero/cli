package subcommands

import "io/ioutil"
import "net/http"
import "os"

const latestDescription = "Check for a newer version."

// Latest prints checks the running version against the latest available.
var Latest = &Subcommand{
	Tag:         "misc",
	Description: whoAmIDescription,
	Handler: func(args []string, paths Paths) {
		var running string
		if args[0] == "" {
			running = "Development Build"
		} else {
			running = "v" + args[0]
		}
		response, err := http.Get("https://licensezero.com/cli-version")
		if err != nil {
			Fail("Could not fetch latest version from licensezero.com.")
		}
		defer response.Body.Close()
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			Fail("Error reading response body.")
		}
		current := string(responseBody)
		os.Stdout.WriteString("Running: " + running + "\n")
		os.Stdout.WriteString("Latest:  " + current + "\n")
		if running == current {
			os.Exit(0)
		} else {
			response, err := http.Get("https://licensezero.com/one-line-install.sh")
			if err != nil {
				os.Exit(1)
			}
			defer response.Body.Close()
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				os.Exit(1)
			}
			os.Stdout.WriteString("Install: " + string(responseBody) + "\n")
			os.Exit(1)
		}
	},
}
