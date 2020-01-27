package data

import "encoding/json"
import "io/ioutil"
import "path"

// LicensorAccount contains the licensor's ID at a specific vendor API.
type LicensorAccount struct {
	LicensorID string `json:"licensorID"`
	API        string `json:"api"`
	Token      string `json:"token"`
}

func licensorPath(home string) string {
	return path.Join(ConfigPath(home), "licensor.json")
}

// ReadLicensorAccounts reads the user's licensor ID and access token from disk.
func ReadLicensorAccounts(home string) ([]LicensorAccount, error) {
	path := licensorPath(home)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var accounts []LicensorAccount
	json.Unmarshal(data, &accounts)
	return accounts, nil
}

// WriteLicensorAccounts writes a licensor ID and access token to disk.
func WriteLicensorAccounts(home string, accounts []LicensorAccount) error {
	data, jsonError := json.Marshal(accounts)
	if jsonError != nil {
		return jsonError
	}
	directoryError := makeConfigDirectory(home)
	if directoryError != nil {
		return directoryError
	}
	return ioutil.WriteFile(licensorPath(home), data, 0644)
}
