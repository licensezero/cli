package data

import "encoding/json"
import "io/ioutil"
import "path"

// Developer describes a developer ID and access token.
type Developer struct {
	Token       string `json:"token"`
	DeveloperID string `json:"developerID"`
}

func developerPath(home string) string {
	return path.Join(ConfigPath(home), "developer.json")
}

// ReadDeveloper reads the user's developer ID and access token from disk.
func ReadDeveloper(home string) (*Developer, error) {
	path := developerPath(home)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var developer Developer
	json.Unmarshal(data, &developer)
	return &developer, nil
}

// WriteDeveloper writes a developer ID and access token to disk.
func WriteDeveloper(home string, developer *Developer) error {
	data, jsonError := json.Marshal(developer)
	if jsonError != nil {
		return jsonError
	}
	directoryError := makeConfigDirectory(home)
	if directoryError != nil {
		return directoryError
	}
	return ioutil.WriteFile(developerPath(home), data, 0644)
}
