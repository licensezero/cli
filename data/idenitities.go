package data

import "encoding/json"
import "io/ioutil"
import "path"

type Identity struct {
	Name         string `json:"name"`
	Jurisdiction string `json:"jurisdiction"`
	EMail        string `json:"email"`
}

func identityPath(home string) string {
	return path.Join(configPath(home), "identity.json")
}

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

func WriteIdentity(home string, identity *Identity) error {
	data, jsonError := json.Marshal(identity)
	if jsonError != nil {
		return jsonError
	}
	directoryError := makeConfigDirectory(home)
	if directoryError != nil {
		return directoryError
	}
	return ioutil.WriteFile(identityPath(home), data, 0744)
}
