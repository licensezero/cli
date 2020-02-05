package cli

import (
	"encoding/json"
	"errors"
	"github.com/yookoala/realpath"
	"io/ioutil"
	"os"
	"path"
)

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
				if alreadyHave(findings, finding) {
					continue
				}
				packageInfo := findPackageInfo(cwd)
				if packageInfo != nil {
					finding.Type = packageInfo.Type
					finding.Name = packageInfo.Name
					finding.Version = packageInfo.Version
					finding.Scope = packageInfo.Scope
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

func findPackageInfo(directoryPath string) *Finding {
	approaches := []func(string) *Finding{
		findNPMPackageInfo,
		// TODO: findPythonPackageInfo,
		// TODO: findMavenPackageInfo,
		// TODO: findComposerPackageInfo,
	}
	for _, approach := range approaches {
		returned := approach(directoryPath)
		if returned != nil {
			return returned
		}
	}
	return nil
}

// localFindings reads project metadata from various files.
func localFindings(directoryPath string) (findings []*Finding, err error) {
	var hadFindings = 0
	var readerFunctions = []func(string) ([]*Finding, error){
		readLicenseZeroJSON,
		// ReadCargoTOML,
	}
	for _, readerFunction := range readerFunctions {
		projects, err := readerFunction(directoryPath)
		if err == nil {
			hadFindings = hadFindings + 1
			findings = projects
		}
	}
	if hadFindings > 1 {
		return nil, errors.New("multiple metadata files")
	}
	return
}

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
		item := Finding{
			API:     offer.API,
			OfferID: offer.OfferID,
			Path:    directoryPath,
			Public:  offer.Public,
		}
		realDirectory, err := realpath.Realpath(directoryPath)
		if err != nil {
			item.Path = realDirectory
		} else {
			item.Path = directoryPath
		}
		findings = append(findings, &item)
	}
	return findings, nil
}
