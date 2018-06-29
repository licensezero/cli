package inventory

import "encoding/json"
import "github.com/yookoala/realpath"
import "io/ioutil"
import "os"
import "path"

type LicenseZeroJSONFile struct {
	Version   string                    `json:"version"`
	Envelopes []ProjectManifestEnvelope `json:"licensezero"`
}

func ReadLicenseZeroFiles(directoryPath string) ([]Project, error) {
	var returned []Project
	entries, err := readAndStatDir(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Project{}, nil
		} else {
			return nil, err
		}
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == "licensezero.json" {
			json_file := path.Join(directoryPath, "licensezero.json")
			data, err := ioutil.ReadFile(json_file)
			if err != nil {
				return nil, err
			}
			var parsed LicenseZeroJSONFile
			json.Unmarshal(data, &parsed)
			for _, envelope := range parsed.Envelopes {
				if alreadyHaveProject(returned, envelope.Manifest.ProjectID) {
					continue
				}
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
			below, err := ReadLicenseZeroFiles(directory)
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
		findPythonPackageInfo,
		findMavenPackageInfo,
	}
	for _, approach := range approaches {
		returned := approach(directoryPath)
		if returned != nil {
			return returned
		}
	}
	return nil
}
