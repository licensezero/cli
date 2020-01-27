package inventory

import "encoding/json"
import "github.com/yookoala/realpath"
import "errors"
import "io/ioutil"
import "os"
import "path"

func recurseLicenseZeroFiles(directoryPath string) ([]DescenderResult, error) {
	var returned []DescenderResult
	entries, err := readAndStatDir(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []DescenderResult{}, nil
		}
		return nil, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == "licensezero.json" {
			results, err := ReadLicenseZeroJSON(directoryPath)
			if err != nil {
				return nil, err
			}
			for _, result := range results {
				if alreadyHave(returned, &result.Offer) {
					continue
				}
				packageInfo := findPackageInfo(directoryPath)
				if packageInfo != nil {
					result.Type = packageInfo.Type
					result.Name = packageInfo.Name
					result.Version = packageInfo.Version
					result.Scope = packageInfo.Scope
				}
				returned = append(returned, result)
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

func findPackageInfo(directoryPath string) *DescenderResult {
	approaches := []func(string) *DescenderResult{
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

// ReadLocalOffers reads offer metadata from various files.
func ReadLocalOffers(directoryPath string) ([]DescenderResult, error) {
	var results []DescenderResult
	var hadResults = 0
	var readerFunctions = []func(string) ([]DescenderResult, error){ReadLicenseZeroJSON, ReadCargoTOML}
	for _, readerFunction := range readerFunctions {
		offers, err := readerFunction(directoryPath)
		if err == nil {
			hadResults = hadResults + 1
			results = offers
		}
	}
	if hadResults > 1 {
		return nil, errors.New("multiple metadata files")
	}
	return results, nil
}

// ReadLicenseZeroJSON reads metadata from licensezero.json.
func ReadLicenseZeroJSON(directoryPath string) ([]DescenderResult, error) {
	var returned []DescenderResult
	jsonFile := path.Join(directoryPath, "licensezero.json")
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	var parsed ArtifactMetadata
	json.Unmarshal(data, &parsed)
	for _, offer := range parsed.Offers {
		offer := DescenderResult{
			Path:  directoryPath,
			Offer: offer,
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
