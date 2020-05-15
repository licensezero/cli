package data

import "encoding/json"
import "io/ioutil"
import "os"
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
		// Attempt to read legacy licensor.json file.
		legacy, legacyErr := readLicensor(home)
		if legacyErr != nil {
			return nil, err
		}
		// If we have a legacy licensor.json file,
		// write a a replacement developer.json file
		// and delete the old licensor.json file.
		upgraded := Developer{
			Token:       legacy.Token,
			DeveloperID: legacy.LicensorID,
		}
		writeErr := WriteDeveloper(home, &upgraded)
		if writeErr == nil {
			os.Remove(licensorPath(home))
		}
		return &upgraded, nil
	}
	var developer Developer
	json.Unmarshal(data, &developer)
	return &developer, nil
}

type legacyLicensor struct {
	Token      string `json:"token"`
	LicensorID string `json:"licensorID"`
}

func licensorPath(home string) string {
	return path.Join(ConfigPath(home), "licensor.json")
}

func readLicensor(home string) (*legacyLicensor, error) {
	path := licensorPath(home)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var licensor legacyLicensor
	json.Unmarshal(data, &licensor)
	return &licensor, nil
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
