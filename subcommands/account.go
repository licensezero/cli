package subcommands

import (
	"licensezero.com/licensezero/user"
)

func selectAccount(broker *string) (*user.Account, string) {
	// Find the relevant account.
	accounts, err := user.ReadAccounts()
	if err != nil {
		return nil, "Error reading accounts: " + err.Error() + "\n"
	}
	if len(accounts) == 0 {
		return nil, "No accounts. Register for an account and run `licensezero token` first.\n"
	}
	base := "https://" + *broker
	if len(accounts) == 1 {
		account := accounts[0]
		if base != "" && account.Server != base {
			return nil, "No account for broker " + *broker + "\n"
		}
		return accounts[0], ""
	}
	if *broker == "" {
		return nil, "You have multiple seller accounts.\nChoose one by passing --broker URL.\n"
	}
	for _, account := range accounts {
		if account.Server == base {
			return account, ""
		}
	}
	return nil, "No account for " + *broker + " found.\n"
}
