package subcommands

import "flag"
import "io/ioutil"

const bugsDescription = "Open the CLI tracker page."

var Bugs = Subcommand{
	Tag:         "misc",
	Description: bugsDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("bugs", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = bugsUsage
		flagSet.Parse(args)
		openURLAndExit("https://github.com/licensezero/cli/issues", doNotOpen)
	},
}

func bugsUsage() {
	usage := bugsDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero bugs\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"do-not-open": doNotOpenLine,
		})
	Fail(usage)
}
