package subcommands

import (
	"errors"
	"flag"
	"io/ioutil"
	"licensezero.com/licensezero/user"
)

const tokenDescription = "Save your API access token."

var tokenUsage = tokenDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero token --seller ID [--broker URL]\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"broker URL": brokerFlagUsage,
		"seller ID":  sellerIDUsage,
		"silent":     silentUsage,
	})

// Token saves licensor IDs and API tokens.
var Token = &Subcommand{
	Tag:         "seller",
	Description: tokenDescription,
	Handler: func(env Environment) int {
		// Parse flags.
		flagSet := flag.NewFlagSet("token", flag.ContinueOnError)
		sellerID := sellerIDFlag(flagSet)
		broker := brokerFlag(flagSet)
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		printUsage := func() {
			env.Stderr.WriteString(tokenUsage)
		}
		flagSet.Usage = printUsage
		err := flagSet.Parse(env.Arguments)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				printUsage()
			}
			return 1
		}
		if *sellerID == "" {
			env.Stderr.WriteString(tokenUsage)
			return 1
		}
		// Prompt for token.
		token, err := env.Stdin.SecretPrompt("Token: ", env.Stdout)
		if err != nil {
			env.Stderr.WriteString(err.Error() + "\n")
			return 1
		}
		base := "https://" + *broker
		// Create account.
		account := user.Account{
			Server:   base,
			SellerID: *sellerID,
			Token:    token,
		}
		// Check existing accounts.
		accounts, err := user.ReadAccounts()
		if err != nil {
			env.Stderr.WriteString("Could not read existing accounts.\n")
			return 1
		}
		for _, existing := range accounts {
			sameAccount := existing.Server == account.Server &&
				existing.SellerID == account.SellerID
			if !sameAccount {
				continue
			}
			if existing.Token == account.Token {
				env.Stderr.WriteString("Already saved.\n")
				return 1
			}
			err = user.DeleteAccount(existing)
			if err != nil {
				env.Stderr.WriteString("Error deleting existing token for " + existing.Server + " ID " + existing.SellerID + ":" + err.Error() + "\n")
				return 1
			}
		}
		err = user.WriteAccount(&account)
		if err != nil {
			env.Stderr.WriteString("Error saving account: " + err.Error() + "\n")
			return 1
		}
		if !*silent {
			env.Stdout.WriteString("Saved.\n")
		}
		return 0
	},
}
