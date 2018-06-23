package inventory

import "encoding/json"
import "io/ioutil"
import "os"
import "path"
import "strings"
import "github.com/yookoala/realpath"

type PackageJSONFile struct {
	Name      string                    `json:"name"`
	Version   string                    `json:"version"`
	Envelopes []ProjectManifestEnvelope `json:"licensezero"`
}

func ReadNPMProjects(packagePath string) ([]Project, error) {
	var returned []Project
	node_modules := path.Join(packagePath, "node_modules")
	entries, err := readAndStatDir(node_modules)
	if err != nil {
		if os.IsNotExist(err) {
			return []Project{}, nil
		} else {
			return nil, err
		}
	}
	processProject := func(directory string, scope *string) error {
		package_json := path.Join(directory, "package.json")
		data, err := ioutil.ReadFile(package_json)
		if err != nil {
			return err
		}
		var parsed PackageJSONFile
		json.Unmarshal(data, &parsed)
		anyNewProjects := false
		for _, envelope := range parsed.Envelopes {
			if alreadyHaveProject(returned, envelope.Manifest.ProjectID) {
				continue
			}
			anyNewProjects = true
			project := Project{
				Type:     "npm",
				Path:     directory,
				Name:     parsed.Name,
				Version:  parsed.Version,
				Envelope: envelope,
			}
			realDirectory, err := realpath.Realpath(directory)
			if err != nil {
				project.Path = realDirectory
			} else {
				project.Path = directory
			}
			if scope != nil {
				project.Scope = *scope
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
		if strings.HasPrefix(name, "@") { // ./node_modules/@scope/package
			scope := name[1:]
			scopePath := path.Join(node_modules, name)
			scopeEntries, err := readAndStatDir(scopePath)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				} else {
					return nil, err
				}
			}
			for _, scopeEntry := range scopeEntries {
				if !scopeEntry.IsDir() {
					continue
				}
				directory := path.Join(node_modules, name, scopeEntry.Name())
				err := processProject(directory, &scope)
				if err != nil {
					return nil, err
				}
			}
		} else { // ./node_modules/package/
			directory := path.Join(node_modules, name)
			err := processProject(directory, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	return returned, nil
}

func alreadyHaveProject(projects []Project, projectID string) bool {
	for _, other := range projects {
		if other.Envelope.Manifest.ProjectID == projectID {
			return true
		}
	}
	return false
}
