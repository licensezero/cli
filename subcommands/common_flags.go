package subcommands

import "flag"

func DoNotOpen(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("do-not-open", false, "Do not open checkout page.")
}

func NoNoncommercial(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("no-noncommercial", false, noNoncommercialLine)
}

func NoReciprocal(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("no-reciprocal", false, noReciprocalLine)
}

func ProjectID(flagSet *flag.FlagSet) *string {
	return flagSet.String("project-id", "", projectIDLine)
}

func Price(flagSet *flag.FlagSet) *uint {
	return flagSet.Uint("price", 0, priceLine)
}

func Relicense(flagSet *flag.FlagSet) *uint {
	return flagSet.Uint("relicense", 0, relicenseLine)
}
