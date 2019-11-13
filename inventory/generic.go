package inventory

import "encoding/json"
import "github.com/yookoala/realpath"
import "io/ioutil"
import "os"
import "path"

// LicenseZeroJSONFile describes the contents of licensezero.json.
type LicenseZeroJSONFile struct {
	Version   string                  `json:"version"`
	Envelopes []OfferManifestEnvelope `json:"licensezero"`
}

func recurseLicenseZeroFiles(directoryPath string) ([]Offer, error) {
	var returned []Offer
	entries, err := readAndStatDir(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Offer{}, nil
		}
		return nil, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == "licensezero.json" {
			offers, err := ReadLicenseZeroJSON(directoryPath)
			if err != nil {
				return nil, err
			}
			for _, offer := range offers {
				if alreadyHaveOffer(returned, offer.Envelope.Manifest.OfferID) {
					continue
				}
				packageInfo := findPackageInfo(directoryPath)
				if packageInfo != nil {
					offer.Type = packageInfo.Type
					offer.Name = packageInfo.Name
					offer.Version = packageInfo.Version
					offer.Scope = packageInfo.Scope
				}
				returned = append(returned, offer)
			}
		} else if entry.IsDir() {
			directory := path.Join(directoryPath, name)
			below, err := recurseLicenseZeroFiles(directory)
			if err != nil {
				return nil, err
			}
			returned = append(returned, below...)
		}
	}
	return returned, nil
}

func findPackageInfo(directoryPath string) *Offer {
	approaches := []func(string) *Offer{
		findNPMPackageInfo,
		findPythonPackageInfo,
		findMavenPackageInfo,
		findComposerPackageInfo,
	}
	for _, approach := range approaches {
		returned := approach(directoryPath)
		if returned != nil {
			return returned
		}
	}
	return nil
}

// ReadLicenseZeroJSON read metadata from licensezero.json.
func ReadLicenseZeroJSON(directoryPath string) ([]Offer, error) {
	var returned []Offer
	jsonFile := path.Join(directoryPath, "licensezero.json")
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	var parsed LicenseZeroJSONFile
	json.Unmarshal(data, &parsed)
	for _, envelope := range parsed.Envelopes {
		offer := Offer{
			Path:     directoryPath,
			Envelope: envelope,
		}
		realDirectory, err := realpath.Realpath(directoryPath)
		if err != nil {
			offer.Path = realDirectory
		} else {
			offer.Path = directoryPath
		}
		returned = append(returned, offer)
	}
	return returned, nil
}
