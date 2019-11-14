package inventory

import "encoding/json"
import "github.com/yookoala/realpath"
import "errors"
import "io/ioutil"
import "os"
import "path"

// LicenseZeroJSONFile describes the contents of licensezero.json.
type LicenseZeroJSONFile struct {
	Version   string                    `json:"version"`
	Envelopes []ProjectManifestEnvelope `json:"licensezero"`
}

func recurseLicenseZeroFiles(directoryPath string) ([]Project, error) {
	var returned []Project
	entries, err := readAndStatDir(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Project{}, nil
		}
		return nil, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == "licensezero.json" {
			projects, err := ReadLicenseZeroJSON(directoryPath)
			if err != nil {
				return nil, err
			}
			for _, project := range projects {
				if alreadyHaveProject(returned, project.Envelope.Manifest.ProjectID) {
					continue
				}
				packageInfo := findPackageInfo(directoryPath)
				if packageInfo != nil {
					project.Type = packageInfo.Type
					project.Name = packageInfo.Name
					project.Version = packageInfo.Version
					project.Scope = packageInfo.Scope
				}
				returned = append(returned, project)
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

func findPackageInfo(directoryPath string) *Project {
	approaches := []func(string) *Project{
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

// ReadLocalProjects reads project metadata from various files.
func ReadLocalProjects(directoryPath string) ([]Project, error) {
	var results []Project
	var hadResults = 0
	var readerFunctions = []func(string) ([]Project, error){ReadLicenseZeroJSON, ReadCargoTOML}
	for _, readerFunction := range readerFunctions {
		projects, err := readerFunction(directoryPath)
		if err == nil {
			hadResults = hadResults + 1
			results = projects
		}
	}
	if hadResults > 1 {
		return nil, errors.New("multiple metadata files")
	}
	return results, nil
}

// ReadLicenseZeroJSON read metadata from licensezero.json.
func ReadLicenseZeroJSON(directoryPath string) ([]Project, error) {
	var returned []Project
	jsonFile := path.Join(directoryPath, "licensezero.json")
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	var parsed LicenseZeroJSONFile
	json.Unmarshal(data, &parsed)
	for _, envelope := range parsed.Envelopes {
		project := Project{
			Path:     directoryPath,
			Envelope: envelope,
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
