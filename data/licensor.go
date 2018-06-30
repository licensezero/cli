package data

import "encoding/json"
import "io/ioutil"
import "path"

// Licensor describes a licensor ID and access token.
type Licensor struct {
	Token      string `json:"token"`
	LicensorID string `json:"licensorID"`
}

func licensorPath(home string) string {
	return path.Join(ConfigPath(home), "licensor.json")
}

// ReadLicensor reads the user's licensor ID and access token from disk.
func ReadLicensor(home string) (*Licensor, error) {
	path := licensorPath(home)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var licensor Licensor
	json.Unmarshal(data, &licensor)
	return &licensor, nil
}

// WriteLicensor writes a licensor ID and access token to disk.
func WriteLicensor(home string, licensor *Licensor) error {
	data, jsonError := json.Marshal(licensor)
	if jsonError != nil {
		return jsonError
	}
	directoryError := makeConfigDirectory(home)
	if directoryError != nil {
		return directoryError
	}
	return ioutil.WriteFile(licensorPath(home), data, 0644)
}
