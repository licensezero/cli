package inventory

import "github.com/yookoala/realpath"
import "fmt"
import "github.com/BurntSushi/toml"
import "os"
import "encoding/json"
import "os/exec"
import "errors"
import "bytes"
import "path"
import "io/ioutil"

func readCargoCrates(packagePath string) ([]Project, error) {
	var returned []Project
	metadata, err := cargoReadMetadata(packagePath)
	if err != nil {
		return nil, err
	}
	for _, packageRecord := range metadata.Packages {
		metadata := packageRecord.Metadata.LicenseZero
		for _, envelope := range metadata.Envelopes {
			projectID := envelope.Manifest.ProjectID
			if alreadyHaveProject(returned, projectID) {
				continue
			}
			project := Project{
				Type:     "cargo",
				Path:     path.Dir(packageRecord.ManifestPath),
				Name:     packageRecord.Name,
				Version:  packageRecord.Version,
				Envelope: envelope,
			}
			returned = append(returned, project)
		}
	}
	return returned, nil
}

// CargoLicenseZeroMap describes metadata for Cargo crates.
type CargoLicenseZeroMap struct {
	Version   string                    `json:"version" toml:"version"`
	Envelopes []ProjectManifestEnvelope `json:"ids" toml:"ids"`
}

type cargoMetadataMap struct {
	LicenseZero CargoLicenseZeroMap `json:"licensezero" toml:"licensezero"`
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
func ReadCargoTOML(directoryPath string) ([]Project, error) {
	var returned []Project
	cargoTOML := path.Join(directoryPath, "Cargo.toml")
	data, err := ioutil.ReadFile(cargoTOML)
	if err != nil {
		return nil, err
	}
	var parsed cargoTOMLData
	if _, err := toml.Decode(string(data), &parsed); err != nil {
		return nil, errors.New("could not parse Cargo.toml")
	}
	fmt.Printf("%+v\n", parsed)
	for _, envelope := range parsed.Package.Metadata.LicenseZero.Envelopes {
		project := Project{
			Path:     directoryPath,
			Envelope: envelope,
		}
		os.Stdout.WriteString(project.Envelope.Manifest.ProjectID)
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
