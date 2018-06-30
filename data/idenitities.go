package data

import "encoding/json"
import "io/ioutil"
import "path"

// Identity represents a licensee identity.
type Identity struct {
	Name         string `json:"name"`
	Jurisdiction string `json:"jurisdiction"`
	EMail        string `json:"email"`
}

func identityPath(home string) string {
	return path.Join(ConfigPath(home), "identity.json")
}

// ReadIdentity reads the user's identity from disk.
func ReadIdentity(home string) (*Identity, error) {
	path := identityPath(home)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var identity Identity
	json.Unmarshal(data, &identity)
	return &identity, nil
}

// WriteIdentity writes a user identity to disk.
func WriteIdentity(home string, identity *Identity) error {
	data, jsonError := json.Marshal(identity)
	if jsonError != nil {
		return jsonError
	}
	directoryError := makeConfigDirectory(home)
	if directoryError != nil {
		return directoryError
	}
	return ioutil.WriteFile(identityPath(home), data, 0644)
}
