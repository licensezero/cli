package data

import "encoding/json"
import "io/ioutil"
import "os"
import "path"

type LicenseEnvelope struct {
	Manifest  LicenseManifest `json:"license"`
	ProjectID string          `json:"projectID"`
	Document  string          `json:"document"`
	PublicKey string          `json:"publicKey"`
	Signature string          `json:"signature"`
}

type LicenseManifest struct {
	Form    string `json:"FORM"`
	Version string `json:"VERSION"`
	Date    string `json:"date"`
	OrderID string `json:"orderID"`
	Project struct {
		ProjectID   string `json:"projectID"`
		Homepage    string `json:"homepage"`
		Description string `json:"description"`
	}
	Licensee struct {
		Name         string `json:"name"`
		Jurisdiction string `json:"jurisdiction"`
		EMail        string `json:"email"`
	}
	Licensor struct {
		Name         string `json:"name"`
		Jurisdiction string `json:"jurisdiction"`
	}
	Price int `json:"price"`
}

func LicensePath(home string, projectID string) string {
	return path.Join(home, "licenses", projectID+".json")
}

func ReadLicenses(home string) ([]LicenseEnvelope, error) {
	directoryPath := path.Join(configPath(home), "licenses")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []LicenseEnvelope{}, nil
		} else {
			return nil, directoryReadError
		}
	}
	var returned []LicenseEnvelope
	for _, entry := range entries {
		name := entry.Name()
		license, err := readLicense(home, name)
		if err != nil {
			return nil, err
		}
		returned = append(returned, *license)
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
