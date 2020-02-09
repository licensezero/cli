package inventory

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/yookoala/realpath"
	"io/ioutil"
	"os/exec"
	"path"
)

func findCargoCrates(packagePath string) (findings []*Finding, err error) {
	metadata, err := runCargoMetadata(packagePath)
	if err != nil {
		return nil, err
	}
	for _, packageRecord := range metadata.Packages {
		license := packageRecord.License
		licenseZeroMetadata := packageRecord.Metadata.LicenseZero
		var artifact Artifact
		artifact, err := mapToArtifact(licenseZeroMetadata)
		if err != nil {
			continue
		}
		for _, offer := range artifact.Offers {
			finding := Finding{
				Type:    "cargo",
				Path:    path.Dir(packageRecord.ManifestPath),
				Name:    packageRecord.Name,
				Version: packageRecord.Version,
			}
			addArtifactOfferToFinding(&offer, &finding)
			if finding.Public == "" && license != "" {
				finding.Public = license
			}
			findings = append(findings, &finding)
		}
	}
	return
}

type cargoMetadataOutput struct {
	Packages []struct {
		Name         string `json:"name"`
		Version      string `json:"version"`
		ManifestPath string `json:"manifest_path"`
		License      string `json:"license"`
		Metadata     struct {
			LicenseZero interface{} `json:"licensezero"`
		} `json:"metadata"`
	} `json:"packages"`
}

func runCargoMetadata(packagePath string) (*cargoMetadataOutput, error) {
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
	Name     string `toml:"name"`
	Version  string `toml:"version"`
	Metadata struct {
		LicenseZero interface{} `toml:"licensezero"`
	} `toml:"metadata"`
}

func readCargoTOML(directoryPath string) (findings []*Finding, err error) {
	cargoTOML := path.Join(directoryPath, "Cargo.toml")
	data, err := ioutil.ReadFile(cargoTOML)
	if err != nil {
		return
	}
	var parsed cargoTOMLData
	if _, err := toml.Decode(string(data), &parsed); err != nil {
		return nil, errors.New("could not parse Cargo.toml")
	}
	artifact, err := mapToArtifact(parsed.Package.Metadata.LicenseZero)
	if err != nil {
		return
	}
	for _, offer := range artifact.Offers {
		finding := Finding{Path: directoryPath}
		addArtifactOfferToFinding(&offer, &finding)
		realDirectory, err := realpath.Realpath(directoryPath)
		if err != nil {
			finding.Path = realDirectory
		} else {
			finding.Path = directoryPath
		}
		findings = append(findings, &finding)
	}
	return
}
