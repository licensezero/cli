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
