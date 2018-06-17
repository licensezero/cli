package subcommands

import "encoding/json"
import "io/ioutil"
import "path"

type Identity struct {
	Name         string `json:"name"`
	Jurisdiction string `json:"jurisdiction"`
	Email        string `json:"email"`
}

func IdentityPath(home string) string {
	return path.Join(home, ".config", "licensezero", "identity.json")
}

func readIdentity(home string) (*Identity, error) {
	path := IdentityPath(home)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var identity Identity
	json.Unmarshal(data, &identity)
	return &identity, nil
}
