package inventory

import (
	"encoding/json"
	"github.com/yookoala/realpath"
	"io/ioutil"
	"os"
	"path"
)

// findLicenseZeroFiles recurses projects directories,
// identifying and processing subdirectories with
// licensezero.json files in them.
//
// licensezero.json files are the primary ways that
// artifacts within projects indicate that users can byu
// licenses through License Zero. However, for packaging
// systems that don't install packages into the project
// working directory, we need other functions to invoke
// commands like `go list` or `bundle list`, and work
// backwards from their output to the user- or system-level
// paths for dependencies.
func findLicenseZeroFiles(cwd string) (findings []*Finding, err error) {
	entries, err := readAndStatDir(cwd)
	if err != nil {
		if os.IsNotExist(err) {
			return findings, nil
		}
		return nil, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == "licensezero.json" {
			fromJSON, err := readLicenseZeroJSON(cwd)
			if err != nil {
				return nil, err
			}
			for _, finding := range fromJSON {
				if alreadyFound(findings, finding) {
					continue
				}
				packageMetadata := readArtifactMetadata(cwd)
				if packageMetadata != nil {
					finding.Type = packageMetadata.Type
					finding.Name = packageMetadata.Name
					finding.Version = packageMetadata.Version
					finding.Scope = packageMetadata.Scope
				}
				findings = append(findings, finding)
			}
		} else if entry.IsDir() {
			directory := path.Join(cwd, name)
			below, err := findLicenseZeroFiles(directory)
			if err != nil {
				return nil, err
			}
			findings = append(findings, below...)
		}
	}
	return
}

// Given the directory of an artifact, readArtifactMetadata
// attempts to read package-system-specific metadata,
// in order to name, version, and other metadata to the
// finding.
func readArtifactMetadata(directoryPath string) *Finding {
	approaches := []func(string) *Finding{
		readNPMPackageMetadata,
		readPythonPackageMetadata,
		readMavenPackageMetadata,
		readComposerPackageMetadata,
	}
	for _, approach := range approaches {
		returned := approach(directoryPath)
		if returned != nil {
			return returned
		}
	}
	return nil
}

// readLicenseZeroJSON reads the licensezero.json file in a
// given artifact subdirectory.
func readLicenseZeroJSON(directoryPath string) (findings []*Finding, err error) {
	jsonFile := path.Join(directoryPath, "licensezero.json")
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	var unstructured interface{}
	json.Unmarshal(data, &unstructured)
	parsed, err := parseArtifact(unstructured)
	for _, offer := range parsed.Offers {
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
	return findings, nil
}
