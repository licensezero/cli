package subcommands

import "flag"

func doNotOpenFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("do-not-open", false, "Do not open checkout page.")
}

func noParityFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("no-parity", false, noReciprocalLine)
}

func noProsperityFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("no-prosperity", false, noNoncommercialLine)
}

func noNoncommercialFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("no-noncommercial", false, noNoncommercialLine)
}

func noReciprocalFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("no-reciprocal", false, noReciprocalLine)
}

func idFlag(flagSet *flag.FlagSet) *string {
	return flagSet.String("id", "", idLine)
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

func noRelicenseFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("no-relicense", false, noRelicenseLine)
}

func silentFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("silent", false, silentLine)
}
