package subcommands

import "flag"

func DoNotOpen(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("do-not-open", false, "Do not open checkout page.")
}
