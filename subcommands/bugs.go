package subcommands

import (
	"errors"
	"flag"
	"io/ioutil"
)

const bugsDescription = "Open the CLI bug tracker page."

var bugsUsage = bugsDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero bugs\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"do-not-open": doNotOpenUsage,
	})

// Bugs opens the CLI tracker bug tracker page.
var Bugs = &Subcommand{
	Tag:         "misc",
	Description: bugsDescription,
	Handler: func(env Environment) int {
		flagSet := flag.NewFlagSet("bugs", flag.ExitOnError)
		doNotOpen := doNotOpenFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		printUsage := func() {
			env.Stderr.WriteString(bugsUsage)
		}
		flagSet.Usage = printUsage
		err := flagSet.Parse(env.Arguments)
		if err != nil {
			if errors.Is(err, flag.ErrHelp) {
				printUsage()
			}
			return 1
		}
		openURL("https://github.com/licensezero/cli/issues", doNotOpen, env.Stdout)
		return 0
	},
}
