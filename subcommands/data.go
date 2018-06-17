package subcommands

import "encoding/json"
import "io/ioutil"
import "path"

type Identity struct {
	Name         string `json:"name"`
	Jurisdiction string `json:"jurisdiction"`
	EMail        string `json:"email"`
}

func identityPath(home string) string {
	return path.Join(home, ".config", "licensezero", "identity.json")
}

func readIdentity(home string) (*Identity, error) {
	path := identityPath(home)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var identity Identity
	json.Unmarshal(data, &identity)
	return &identity, nil
}

func writeIdentity(home string, identity *Identity) error {
	data, err := json.Marshal(identity)
	if err != nil {
		return err
	}
	path := identityPath(home)
	return ioutil.WriteFile(path, data, 0644)
}
