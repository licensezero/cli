package subcommands

import "flag"

func doNotOpenFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("do-not-open", false, doNotOpenLine)
}

func noncommercialFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("noncommercial", false, noncommercialLine)
}

func openFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("open", false, openLine)
}

func sellerIDFlag(flagSet *flag.FlagSet) *string {
	return flagSet.String("seller", "", sellerIDLine)
}

const defaultBroker = "broker.licensezero.com"

var brokerFlagUsage = "Broker server name [default: " + defaultBroker + "]."

func brokerFlag(flagSet *flag.FlagSet) *string {
	return flagSet.String("broker", defaultBroker, brokerFlagUsage)
}

func offerIDFlag(flagSet *flag.FlagSet) *string {
	return flagSet.String("offer", "", offerIDLine)
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
