package user

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Account contains information about a seller account.
type Account struct {
	Server   string `json:"server"`
	SellerID string `json:"sellerID"`
	Token    string `json:"token"`
}

// ReadAccounts reads all the broker server accounts saved
// for the CLI user.
func ReadAccounts() (accounts []*Account, err error) {
	configPath, err := ConfigPath()
	if err != nil {
		return
	}
	directoryPath := path.Join(configPath, "accounts")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return
		}
		return nil, directoryReadError
	}
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(configPath, "accounts", name)
		account, err := readAccount(filePath)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return
}

func readAccount(filePath string) (account *Account, err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &account)
	return
}

func WriteAccount(account *Account) (err error) {
	configPath, err := ConfigPath()
	if err != nil {
		return
	}
	directoryPath := path.Join(configPath, "accounts")
	err = os.MkdirAll(directoryPath, 0700)
	if err != nil {
		return
	}
	filePath := path.Join(directoryPath, account.Token+".json")
	data, err := json.Marshal(account)
	if err != nil {
		return
	}
	return ioutil.WriteFile(filePath, data, 0644)
}

func DeleteAccount(account *Account) (err error) {
	configPath, err := ConfigPath()
	if err != nil {
		return
	}
	filePath := path.Join(configPath, "accounts", account.Token+".json")
	err = os.Remove(filePath)
	return
}
