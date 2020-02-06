package user

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

// Identity contains information about the CLI user.
type Identity struct {
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
}

// ReadIdentity reads the CLI user's identify.
func ReadIdentity(configPath string) (identity *Identity, err error) {
	filePath := path.Join(configPath, "identity.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &identity)
	return
}
