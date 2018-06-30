package subcommands

import "flag"

func doNotOpenFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("do-not-open", false, "Do not open checkout page.")
}

func noNoncommercialFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("no-noncommercial", false, noNoncommercialLine)
}

func noReciprocalFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("no-reciprocal", false, noReciprocalLine)
}

func projectIDFlag(flagSet *flag.FlagSet) *string {
	return flagSet.String("project", "", projectIDLine)
}

func priceFlag(flagSet *flag.FlagSet) *uint {
	return flagSet.Uint("price", 0, priceLine)
}

func relicenseFlag(flagSet *flag.FlagSet) *uint {
	return flagSet.Uint("relicense", 0, relicenseLine)
}

func silentFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("silent", false, silentLine)
}
