package data

import "encoding/json"
import "io/ioutil"
import "path"

type Licensor struct {
	Token      string `json:"token"`
	LicensorID string `json:"licensorID"`
}

func licensorPath(home string) string {
	return path.Join(configPath(home), "licensor.json")
}

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

func WriteLicensor(home string, licensor *Licensor) error {
	data, jsonError := json.Marshal(licensor)
	if jsonError != nil {
		return jsonError
	}
	directoryError := makeConfigDirectory(home)
	if directoryError != nil {
		return directoryError
	}
	return ioutil.WriteFile(licensorPath(home), data, 0744)
}
