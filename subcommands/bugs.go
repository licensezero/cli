package subcommands

import "flag"
import "os"

const bugsDescription = "Access the bug tracker for the application."

var Bugs = Subcommand{
	Description: bugsDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("bugs", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
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
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
