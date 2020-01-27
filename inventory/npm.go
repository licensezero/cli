package inventory

import "encoding/json"
import "github.com/yookoala/realpath"
import "io/ioutil"
import "os"
import "path"
import "strings"

type packageJSONFile struct {
	Name     string           `json:"name"`
	Version  string           `json:"version"`
	Artifact ArtifactMetadata `json:"licensezero"`
}

func readNPMProjects(packagePath string) ([]DescenderResult, error) {
	var returned []DescenderResult
	nodeModules := path.Join(packagePath, "node_modules")
	entries, err := readAndStatDir(nodeModules)
	if err != nil {
		if os.IsNotExist(err) {
			return []DescenderResult{}, nil
		}
		return nil, err
	}
	processProject := func(directory string, scope *string) error {
		anyNewProjects := false
		parsed, err := readPackageJSON(directory)
		if err != nil {
			return err
		}
		for _, offer := range parsed.Artifact.Offers {
			if alreadyHave(returned, &offer) {
				continue
			}
			anyNewProjects = true
			result := DescenderResult{
				Type:    "npm",
				Path:    directory,
				Name:    parsed.Name,
				Version: parsed.Version,
				Offer:   offer,
			}
			realDirectory, err := realpath.Realpath(directory)
			if err != nil {
				result.Path = realDirectory
			} else {
				result.Path = directory
			}
			if scope != nil {
				result.Name = *scope + result.Name
			}
			returned = append(returned, result)
		}
		if anyNewProjects {
			below, recursionError := readNPMProjects(directory)
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
			scopePath := path.Join(nodeModules, name)
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
				directory := path.Join(nodeModules, name, scopeEntry.Name())
				err := processProject(directory, &scope)
				if err != nil {
					return nil, err
				}
			}
		} else { // ./node_modules/package/
			directory := path.Join(nodeModules, name)
			err := processProject(directory, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	return returned, nil
}

func readPackageJSON(directory string) (*packageJSONFile, error) {
	packageJSON := path.Join(directory, "package.json")
	data, err := ioutil.ReadFile(packageJSON)
	if err != nil {
		return nil, err
	}
	var parsed packageJSONFile
	json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func findNPMPackageInfo(directoryPath string) *DescenderResult {
	packageJSON := path.Join(directoryPath, "package.json")
	data, err := ioutil.ReadFile(packageJSON)
	if err != nil {
		return nil
	}
	var parsed struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return nil
	}
	rawName := parsed.Name
	var name, scope string
	// If the name looks like @scope/name, parse it.
	if strings.HasPrefix(rawName, "@") && strings.Index(rawName, "/") != -1 {
		index := strings.Index(rawName, "/")
		scope = rawName[1 : index-1]
		name = rawName[index:]
	} else {
		name = rawName
	}
	return &DescenderResult{
		Type:    "npm",
		Name:    name,
		Scope:   scope,
		Version: parsed.Version,
	}
}
