package inventory

import "encoding/json"
import "os/exec"
import "bytes"
import "path"

func readCargoCrates(packagePath string) ([]Offer, error) {
	var returned []Offer
	metadata, err := cargoReadMetadata(packagePath)
	if err != nil {
		return nil, err
	}
	for _, packageRecord := range metadata.Packages {
		metadata := packageRecord.Metadata.LicenseZero
		for _, envelope := range metadata.Envelopes {
			offerID := envelope.Manifest.OfferID
			if alreadyHaveOffer(returned, offerID) {
				continue
			}
			offer := Offer{
				Type:     "cargo",
				Path:     path.Dir(packageRecord.ManifestPath),
				Name:     packageRecord.Name,
				Version:  packageRecord.Version,
				Envelope: envelope,
			}
			returned = append(returned, offer)
		}
	}
	return returned, nil
}

type CargoLicenseZeroMap struct {
	Version   string                  `json:"version"`
	Envelopes []OfferManifestEnvelope `json:"ids"`
}

type cargoMetadataMap struct {
	LicenseZero CargoLicenseZeroMap `json:"licensezero"`
}

type cargoMetadataPackage struct {
	Name         string           `json:"name"`
	Version      string           `json:"version"`
	ManifestPath string           `json:"manifest_path"`
	Metadata     cargoMetadataMap `json:"metadata"`
}

type cargoMetadataOutput struct {
	Packages []cargoMetadataPackage `json:"packages"`
}

func cargoReadMetadata(packagePath string) (*cargoMetadataOutput, error) {
	list := exec.Command("cargo", "metadata", "--format-version", "1")
	list.Dir = packagePath
	var stdout bytes.Buffer
	list.Stdout = &stdout
	err := list.Run()
	if err != nil {
		return nil, err
	}
	var parsed cargoMetadataOutput
	json.Unmarshal(stdout.Bytes(), &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
