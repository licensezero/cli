package inventory

import "encoding/json"
import "github.com/yookoala/realpath"
import "github.com/mitchellh/mapstructure"
import "errors"
import "io/ioutil"
import "os"
import "path"

// Version2LicenseZeroMetadata describes the contents of a version 2 metadata.
type Version2LicenseZeroMetadata []Version2Envelope

// Version1LicenseZeroMetadata describes the contents of a version 1 metadata.
type Version1LicenseZeroMetadata struct {
	Version   string             `json:"version" mapstructure:"version"`
	Envelopes []Version1Envelope `json:"licensezero" mapstructure:"licensezero"`
}

func (json Version1LicenseZeroMetadata) offers() []Offer {
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

// ReadLicenseZeroJSON reads metadata from licensezero.json.
func ReadLicenseZeroJSON(directoryPath string) ([]Offer, error) {
	realDirectory, err := realpath.Realpath(directoryPath)
	if err != nil {
		directoryPath = realDirectory
	}
	jsonFile := path.Join(directoryPath, "licensezero.json")
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	var unstructured interface{}
	json.Unmarshal(data, &unstructured)
	returned, err := interpretLicenseZeroMetadata(unstructured)
	if err != nil {
		return nil, err
	}
	for _, offer := range returned {
		offer.Artifact.Path = directoryPath
	}
	return returned, nil
}

func interpretLicenseZeroMetadata(unstructured interface{}) ([]Offer, error) {
	switch value := unstructured.(type) {
	case []interface{}:
		return parseArrayMetadata(value)
	case map[string]interface{}:
		return parseObjectMetadata(value)
	default:
		return nil, errors.New("could not parse licensezero.json")
	}
}

func parseArrayMetadata(parsed []interface{}) ([]Offer, error) {
	var returned []Offer
	// Iterate elements of the JSON Array.
	for _, entry := range parsed {
		object, matched := entry.(map[string]interface{})
		if !matched {
			return nil, errors.New("invalid entry")
		}
		schema, matched := object["schema"].(string)
		if !matched || schema == "" {
			// Version 1 envelope.
			var envelope Version1Envelope
			err := mapstructure.Decode(entry, &envelope)
			if err != nil {
				return nil, errors.New("invalid entry")
			}
			var offer = envelope.offer()
			returned = append(returned, offer)
		} else if schema == "2.0.0" {
			// Version 2 envelope.
			var envelope Version2Envelope
			err := mapstructure.Decode(entry, &envelope)
			if err != nil {
				return nil, errors.New("invalid entry")
			}
			var offer = envelope.offer()
			returned = append(returned, offer)
		} else {
			// Unknown version schema.
			// TODO: Show hint to run `latest` on encountering unknown schema.
			return nil, errors.New("unkown version schema")
		}
	}
	return returned, nil
}

func parseObjectMetadata(unstructured interface{}) ([]Offer, error) {
	var metadata Version1LicenseZeroMetadata
	err := mapstructure.Decode(unstructured, &metadata)
	if err != nil {
		return nil, errors.New("could not parse licensezero.json")
	}
	var returned []Offer
	for _, envelope := range metadata.Envelopes {
		offer := Offer{
			License: LicenseData{
				Terms:   envelope.Manifest.Terms,
				Version: envelope.Manifest.Version,
			},
			OfferID:  envelope.Manifest.ProjectID,
			Envelope: envelope,
		}
		returned = append(returned, offer)
	}
	return returned, nil
}
