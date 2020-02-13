package subcommands

import "flag"

const doNotOpenUsage = "Do not open pages in web browser."

func doNotOpenFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("do-not-open", false, doNotOpenUsage)
}

const noncommercialUsage = "Noncommercial project. Ignore noncommercial licenses."

func noncommercialFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("noncommercial", false, noncommercialUsage)
}

const openUsage = "Open software project. Ignore reciprocal licenses."

func openFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("open", false, openUsage)
}

const sellerIDUsage = "Seller ID (UUID)."

func sellerIDFlag(flagSet *flag.FlagSet) *string {
	return flagSet.String("seller", "", sellerIDUsage)
}

const defaultBroker = "broker.licensezero.com"

var brokerFlagUsage = "Broker server name [default: " + defaultBroker + "]."

func brokerFlag(flagSet *flag.FlagSet) *string {
	return flagSet.String("broker", defaultBroker, brokerFlagUsage)
}

const offerIDUsage = "Offer ID (UUID)."

func offerIDFlag(flagSet *flag.FlagSet) *string {
	return flagSet.String("offer", "", offerIDUsage)
}

const priceUsage = "Private license price, in US cents."

func priceFlag(flagSet *flag.FlagSet) *uint {
	return flagSet.Uint("price", 0, priceUsage)
}

const silentUsage = "Suppress output about success."

func silentFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("silent", false, silentUsage)
}
