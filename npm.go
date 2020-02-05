package cli

import (
	"encoding/json"
	"github.com/yookoala/realpath"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func findNPMPackages(packagePath string) (findings []*Finding, err error) {
	nodeModules := path.Join(packagePath, "node_modules")
	entries, err := readAndStatDir(nodeModules)
	if err != nil {
		if os.IsNotExist(err) {
			return findings, nil
		}
		return
	}

	processFinding := func(directory string, scope *string) (err error) {
		anyNewFindings := false
		packageJSON, err := readPackageJSON(directory)
		if err != nil {
			return
		}
		artifact, err := parseArtifact(packageJSON.LicenseZero)
		if err != nil {
			return
		}
		for _, offer := range artifact.Offers {
			finding := Finding{
				Type:    "npm",
				Path:    directory,
				Name:    packageJSON.Name,
				Version: packageJSON.Version,
			}
			addArtifactOfferToFinding(&offer, &finding)
			if alreadyHave(findings, &finding) {
				continue
			}
			anyNewFindings = true
			realDirectory, err := realpath.Realpath(directory)
			if err != nil {
				finding.Path = realDirectory
			} else {
				finding.Path = directory
			}
			if scope != nil {
				finding.Scope = *scope
			}
			findings = append(findings, &finding)
		}
		if anyNewFindings {
			below, recursionError := findNPMPackages(directory)
			if recursionError != nil {
				return recursionError
			}
			findings = append(findings, below...)
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
				}
				return nil, err
			}
			for _, scopeEntry := range scopeEntries {
				if !scopeEntry.IsDir() {
					continue
				}
				directory := path.Join(nodeModules, name, scopeEntry.Name())
				err = processFinding(directory, &scope)
				if err != nil {
					return nil, err
				}
			}
		} else { // ./node_modules/package/
			directory := path.Join(nodeModules, name)
			err = processFinding(directory, nil)
			if err != nil {
				return
			}
		}
	}
	return
}

type packageJSONFile struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	LicenseZero interface{} `json:"licensezero"`
}

func readPackageJSON(directory string) (read *packageJSONFile, err error) {
	packageJSON := path.Join(directory, "package.json")
	data, err := ioutil.ReadFile(packageJSON)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &read)
	if err != nil {
		return
	}
	return
}

func readNPMPackageInfo(packagePath string) *Finding {
	packageJSON := path.Join(packagePath, "package.json")
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
	return &Finding{
		Type:    "npm",
		Name:    name,
		Scope:   scope,
		Version: parsed.Version,
	}
}
