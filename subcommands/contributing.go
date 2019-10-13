package subcommands

import "flag"
import "github.com/licensezero/cli/contributing"
import "io/ioutil"
import "os"
import "strings"

const contributingDescription = "Add licensing information to CONTRIBUTING."

// Contributing appends licensing information to CONTRIBUTING.
var Contributing = &Subcommand{
	Tag:         "seller",
	Description: contributingDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("contributing", flag.ExitOnError)
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			usage := contributingDescription + "\n\n" +
				"Usage:\n" +
				"  licensezero contributing\n\n" +
				"Options:\n" +
				flagsList(map[string]string{
					"silent": silentLine,
				})
			Fail(usage)
		}
		flagSet.Parse(args)
		// Append to CONTRIBUTING.
		contributingName, data, err := contributing.ReadCONTRIBUTING(paths.CWD)
		var toWrite string
		if err != nil {
			Fail("Error: " + err.Error())
		} else {
			toWrite = string(data)
			if strings.HasSuffix(toWrite, "\n") {
				toWrite = toWrite + "\n"
			} else {
				toWrite = toWrite + "\n\n"
			}
		}
		toWrite = toWrite +
			"# Licensing\n\n" +
			"If you submit a pull request, please be prepared to license " +
			"your contributions under the terms of the " +
			"[Charity Public License](https://licensezero.com/licenses/charity), " +
			"a modern evolution of licenses like MIT and the two-clause BSD license.\n"
		err = ioutil.WriteFile(contributingName, []byte(toWrite), 0644)
		if err != nil {
			Fail("Error writing CONTRIBUTING")
		}
		if !*silent {
			os.Stdout.WriteString("Appended terms to CONTRIBUTING.\n")
		}
		os.Exit(0)
	},
}
