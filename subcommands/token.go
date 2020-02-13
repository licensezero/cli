package subcommands

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/user"
	"net/http"
)

const tokenDescription = "Save your API access token."

var tokenUsage = tokenDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero token --seller ID [--broker URL]\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"broker URL": brokerFlagUsage,
		"seller ID":  sellerIDLine,
		"silent":     silentLine,
	})

// Token saves licensor IDs and API tokens.
var Token = &Subcommand{
	Tag:         "seller",
	Description: tokenDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		// Parse flags.
		flagSet := flag.NewFlagSet("token", flag.ContinueOnError)
		sellerID := sellerIDFlag(flagSet)
		broker := brokerFlag(flagSet)
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			stderr.WriteString(tokenUsage)
		}
		err := flagSet.Parse(args)
		if err != nil {
			if !errors.Is(err, flag.ErrHelp) {
				stderr.WriteString(err.Error() + "\n")
			}
			return 1
		}
		if *sellerID == "" {
			stderr.WriteString(tokenUsage)
			return 1
		}
		// Prompt for token.
		token, err := stdin.SecretPrompt("Token: ", stdout)
		if err != nil {
			stderr.WriteString(err.Error() + "\n")
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
			stderr.WriteString("Could not read existing accounts.\n")
			return 1
		}
		for _, existing := range accounts {
			sameAccount := existing.Server == account.Server &&
				existing.SellerID == account.SellerID
			if !sameAccount {
				continue
			}
			if existing.Token == account.Token {
				stderr.WriteString("Already saved.\n")
				return 1
			}
			err = user.DeleteAccount(existing)
			if err != nil {
				stderr.WriteString("Error deleting existing token for " + existing.Server + " ID " + existing.SellerID + ":" + err.Error() + "\n")
				return 1
			}
		}
		err = user.WriteAccount(&account)
		if err != nil {
			stderr.WriteString("Error saving account: " + err.Error() + "\n")
			return 1
		}
		if !*silent {
			stdout.WriteString("Saved.\n")
		}
		return 0
	},
}
