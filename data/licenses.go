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
	return path.Join(LicensesPath(home), projectID+".json")
}

func LicensesPath(home string) string {
	return path.Join(configPath(home), "licenses")
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
		filePath := path.Join(home, "licenses", name)
		license, err := ReadLicense(filePath)
		if err != nil {
			return nil, err
		}
		returned = append(returned, *license)
	}
	return returned, nil
}

func ReadLicense(filePath string) (*LicenseEnvelope, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var license LicenseEnvelope
	json.Unmarshal(data, &license)
	return &license, nil
}

func WriteLicense(home string, license *LicenseEnvelope) error {
	json, err := json.Marshal(license)
	if err != nil {
		return err
	}
	err = os.MkdirAll(LicensesPath(home), 0744)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(LicensePath(home, license.ProjectID), json, 0744)
}
