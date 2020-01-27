package inventory

import "github.com/yookoala/realpath"
import "github.com/BurntSushi/toml"
import "encoding/json"
import "os/exec"
import "errors"
import "bytes"
import "path"
import "io/ioutil"

func readCargoCrates(packagePath string) ([]DescenderResult, error) {
	var returned []DescenderResult
	metadata, err := cargoReadMetadata(packagePath)
	if err != nil {
		return nil, err
	}
	for _, packageRecord := range metadata.Packages {
		metadata := packageRecord.Metadata.LicenseZero
		for _, offer := range metadata.Offers {
			if alreadyHave(returned, &offer) {
				continue
			}
			project := DescenderResult{
				Type:    "cargo",
				Path:    path.Dir(packageRecord.ManifestPath),
				Name:    packageRecord.Name,
				Version: packageRecord.Version,
				Offer:   offer,
			}
			returned = append(returned, project)
		}
	}
	return returned, nil
}

type cargoMetadataMap struct {
	LicenseZero ArtifactMetadata `json:"licensezero" toml:"licensezero"`
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
func ReadCargoTOML(directoryPath string) ([]DescenderResult, error) {
	var returned []DescenderResult
	cargoTOML := path.Join(directoryPath, "Cargo.toml")
	data, err := ioutil.ReadFile(cargoTOML)
	if err != nil {
		return nil, err
	}
	var parsed cargoTOMLData
	if _, err := toml.Decode(string(data), &parsed); err != nil {
		return nil, errors.New("could not parse Cargo.toml")
	}
	for _, offer := range parsed.Package.Metadata.LicenseZero.Offers {
		project := DescenderResult{
			Path:  directoryPath,
			Offer: offer,
		}
		realDirectory, err := realpath.Realpath(directoryPath)
		if err != nil {
			project.Path = realDirectory
		} else {
			project.Path = directoryPath
		}
		returned = append(returned, project)
	}
	return returned, nil
}
