package inventory

import "encoding/json"
import "github.com/yookoala/realpath"
import "errors"
import "io/ioutil"
import "os"
import "path"

// Version1LicenseZeroJSONFile describes the contents of a version 1 licensezero.json file.
type Version1LicenseZeroJSONFile struct {
	Version   string             `json:"version"`
	Envelopes []Version1Envelope `json:"licensezero"`
}

func (json Version1LicenseZeroJSONFile) offers() []Offer {
	var returned []Offer
	for _, envelope := range json.Envelopes {
		returned = append(returned, Offer{
			OfferID: envelope.Manifest.ProjectID,
			License: LicenseData{
				Terms:   envelope.Manifest.Terms,
				Version: envelope.Manifest.Version,
			},
		})
	}
	return returned
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
				if alreadyHaveOffer(returned, offer.OfferID) {
					continue
				}
				packageInfo := findPackageInfo(directoryPath)
				if packageInfo != nil {
					offer.Artifact = packageInfo.Artifact
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

// ReadPackageOffers reads offer metadata from various files.
func ReadPackageOffers(directoryPath string) ([]Offer, error) {
	var returned []Offer
	var hadOffers = 0
	var readerFunctions = []func(string) ([]Offer, error){ReadLicenseZeroJSON, ReadCargoTOML}
	for _, readerFunction := range readerFunctions {
		offers, err := readerFunction(directoryPath)
		if err == nil {
			hadOffers = hadOffers + 1
			returned = offers
		}
	}
	if hadOffers > 1 {
		return nil, errors.New("multiple metadata files")
	}
	return returned, nil
}

// ReadLicenseZeroJSON read metadata from licensezero.json.
func ReadLicenseZeroJSON(directoryPath string) ([]Offer, error) {
	var returned []Offer
	jsonFile := path.Join(directoryPath, "licensezero.json")
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	var parsed Version1LicenseZeroJSONFile
	json.Unmarshal(data, &parsed)
	for _, envelope := range parsed.Envelopes {
		offer := Offer{
			License: LicenseData{
				Terms:   envelope.Manifest.Terms,
				Version: envelope.Manifest.Version,
			},
		}
		realDirectory, err := realpath.Realpath(directoryPath)
		if err != nil {
			offer.Artifact.Path = realDirectory
		} else {
			offer.Artifact.Path = directoryPath
		}
		returned = append(returned, offer)
	}
	return returned, nil
}
