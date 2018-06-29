package inventory

import "bytes"
import "encoding/json"
import "github.com/yookoala/realpath"
import "io/ioutil"
import "os"
import "os/exec"
import "path"
import "strings"

type LicenseZeroJSONFile struct {
	Version   string                    `json:"version"`
	Envelopes []ProjectManifestEnvelope `json:"licensezero"`
}

// TODO: Consider reading setup.py --url and checking against homepage for Python.

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
	processProject := func(directory string) error {
		json_file := path.Join(directory, "licensezero.json")
		data, err := ioutil.ReadFile(json_file)
		if err != nil {
			return err
		}
		var parsed LicenseZeroJSONFile
		json.Unmarshal(data, &parsed)
		anyNewProjects := false
		for _, envelope := range parsed.Envelopes {
			if alreadyHaveProject(returned, envelope.Manifest.ProjectID) {
				continue
			}
			anyNewProjects = true
			project := Project{
				Path:     directory,
				Envelope: envelope,
			}
			realDirectory, err := realpath.Realpath(directory)
			if err != nil {
				project.Path = realDirectory
			} else {
				project.Path = directory
			}
			packageInfo := findPackageInfo(directory)
			if packageInfo != nil {
				project.Type = packageInfo.Type
				project.Name = packageInfo.Name
				project.Version = packageInfo.Version
				project.Scope = packageInfo.Scope
			}
			returned = append(returned, project)
		}
		if anyNewProjects {
			below, recursionError := ReadNPMProjects(directory)
			if recursionError != nil {
				return recursionError
			}
			returned = append(returned, below...)
		}
		return nil
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		directory := path.Join(directoryPath, name)
		err := processProject(directory)
		if err != nil {
			return nil, err
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

func findPythonPackageInfo(directoryPath string) *Project {
	setup := path.Join(directoryPath, "setup.py")
	_, err := os.Stat(setup)
	if err != nil {
		return nil
	}
	command := exec.Command("python", "setup.py", "--name", "--version")
	var stdout bytes.Buffer
	command.Stdout = &stdout
	err = command.Run()
	if err != nil {
		return nil
	}
	output := string(stdout.Bytes())
	lines := strings.Split(output, "\n")
	name := strings.TrimSpace(lines[0])
	version := strings.TrimSpace(lines[1])
	return &Project{
		Type:    "python",
		Name:    name,
		Version: version,
	}
}
