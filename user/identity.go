package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"path"
)

// Identity contains information about the CLI user.
type Identity struct {
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
}

var ErrEMailMismatch = errors.New("e-mail mismatch")

var ErrJurisdictionMismatch = errors.New("jurisdiction mismatch")

var ErrNameMismatch = errors.New("name mismatch")

func (identity *Identity) ValidateReceipt(receipt *api.Receipt) (errors []error) {
	buyer := receipt.License.Values.Buyer
	if buyer.Name != identity.Name {
		errors = append(errors, ErrNameMismatch)
	}
	if buyer.EMail != identity.EMail {
		errors = append(errors, ErrEMailMismatch)
	}
	if buyer.Jurisdiction != identity.Jurisdiction {
		errors = append(errors, ErrJurisdictionMismatch)
	}
	return
}

// ReadIdentity reads the CLI user's identify.
func ReadIdentity() (identity *Identity, err error) {
	configPath, err := ConfigPath()
	if err != nil {
		return
	}
	data, err := ioutil.ReadFile(identityPath(configPath))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &identity)
	return
}

// WriteIdentity writes the CLI user's identity data to disk.
func WriteIdentity(i *Identity) (err error) {
	data, jsonError := json.Marshal(i)
	if jsonError != nil {
		return jsonError
	}
	configPath, err := ConfigPath()
	if err != nil {
		return
	}
	directoryError := makeConfigDirectory()
	if directoryError != nil {
		return directoryError
	}
	return ioutil.WriteFile(identityPath(configPath), data, 0644)
}

func identityPath(configPath string) string {
	return path.Join(configPath, "identity.json")
}
