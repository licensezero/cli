package inventory

import "github.com/yookoala/realpath"
import "github.com/BurntSushi/toml"
import "encoding/json"
import "os/exec"
import "errors"
import "bytes"
import "path"
import "io/ioutil"

func readCargoCrates(packagePath string) ([]Offer, error) {
	var returned []Offer
	metadata, err := cargoReadMetadata(packagePath)
	if err != nil {
		return nil, err
	}
	for _, packageRecord := range metadata.Packages {
		metadata := packageRecord.Metadata.LicenseZero
		offers, err := interpretLicenseZeroMetadata(metadata)
		if err != nil {
			return nil, err
		}
		for _, offer := range offers {
			offer.Artifact = ArtifactData{
				Type:    "cargo",
				Path:    path.Dir(packageRecord.ManifestPath),
				Name:    packageRecord.Name,
				Version: packageRecord.Version,
			}
			returned = append(returned, offer)
		}
	}
	return returned, nil
}

type cargoMetadataMap struct {
	LicenseZero interface{} `json:"licensezero" toml:"licensezero"`
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

type cargoTOMLData struct {
	Package cargoTOMLPackage `toml:"package"`
}

type cargoTOMLPackage struct {
	Name     string           `toml:"name"`
	Version  string           `toml:"version"`
	Metadata cargoMetadataMap `toml:"metadata"`
}

// ReadCargoTOML reads metadata from Cargo.toml.
func ReadCargoTOML(directoryPath string) ([]Offer, error) {
	var returned []Offer
	cargoTOML := path.Join(directoryPath, "Cargo.toml")
	data, err := ioutil.ReadFile(cargoTOML)
	if err != nil {
		return nil, err
	}
	var parsed cargoTOMLData
	if _, err := toml.Decode(string(data), &parsed); err != nil {
		return nil, errors.New("could not parse Cargo.toml")
	}
	offers, err := interpretLicenseZeroMetadata(parsed.Package.Metadata.LicenseZero)
	if err != nil {
		return nil, err
	}
	realDirectory, err := realpath.Realpath(directoryPath)
	if err != nil {
		directoryPath = realDirectory
	}
	for _, offer := range offers {
		offer.Artifact = ArtifactData{
			Type:    "cargo",
			Name:    parsed.Package.Name,
			Version: parsed.Package.Version,
			Path:    directoryPath,
		}
		returned = append(returned, offer)
	}
	return returned, nil
}
