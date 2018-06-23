package data

import "encoding/json"
import "io/ioutil"
import "os"
import "path"

type LicenseEnvelope struct {
	Manifest LicenseManifest `json:"license"`
}

type LicenseManifest struct {
	ProjectID string `json:"projectID"`
}

func LicensePath(home string, projectID string) string {
	return path.Join(home, "licenses", projectID+".json")
}

func ReadLicenses(home string) ([]LicenseManifest, error) {
	directoryPath := path.Join(configPath(home), "licenses")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []LicenseManifest{}, nil
		} else {
			return nil, directoryReadError
		}
	}
	var returned []LicenseManifest
	for _, entry := range entries {
		name := entry.Name()
		license, err := readLicense(home, name)
		if err != nil {
			return nil, err
		}
		returned = append(returned, license.Manifest)
	}
	return returned, nil
}

func readLicense(home string, file string) (*LicenseEnvelope, error) {
	filePath := path.Join(home, "licenses", file)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var license LicenseEnvelope
	json.Unmarshal(data, &license)
	return &license, nil
}
