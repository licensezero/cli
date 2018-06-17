package subcommands

import "encoding/json"
import "io/ioutil"
import "os"
import "path"

type Identity struct {
	Name         string `json:"name"`
	Jurisdiction string `json:"jurisdiction"`
	EMail        string `json:"email"`
}

func configPath(home string) string {
	return path.Join(home, ".config", "licensezero")
}

func identityPath(home string) string {
	return path.Join(configPath(home), "identity.json")
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

func makeConfigDirectory(home string) error {
	path := configPath(home)
	return os.MkdirAll(path, 0744)
}

type Licensor struct {
	Token      string `json:"token"`
	LicensorID string `json:"licensorID"`
}

func licensorPath(home string) string {
	return path.Join(configPath(home), "licensor.json")
}

func readLicensor(home string) (*Licensor, error) {
	path := licensorPath(home)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var licensor Licensor
	json.Unmarshal(data, &licensor)
	return &licensor, nil
}

func writeLicensor(home string, licensor *Licensor) error {
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

type WaiverEnvelope struct {
	Manifest WaiverManifest `json:"manifest"`
}

type WaiverManifest struct {
	ProjectID    string `json:"projectID"`
	Date         string
	Term         string
	Beneficiary  string
	Jurisdiction string
	EMail        string
}

func WaiverPath(home string, projectID string) string {
	return path.Join(home, "waivers", projectID+".json")
}

func ReadWaivers(home string) ([]WaiverManifest, error) {
	directoryPath := path.Join(configPath(home), "waivers")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []WaiverManifest{}, nil
		} else {
			return nil, directoryReadError
		}
	}
	var returned []WaiverManifest
	for _, entry := range entries {
		name := entry.Name()
		waiver, err := readWaiver(home, name)
		if err != nil {
			return nil, err
		}
		if Unexpired(&waiver.Manifest) {
			returned = append(returned, waiver.Manifest)
		}
	}
	return returned, nil
}

func Unexpired(waiver *WaiverManifest) bool {
	// TODO
	return true
}

func readWaiver(home string, file string) (*WaiverEnvelope, error) {
	filePath := path.Join(home, "waivers", file)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var waiver WaiverEnvelope
	json.Unmarshal(data, &waiver)
	return &waiver, nil
}
