package inventory

import "encoding/json"
import "github.com/yookoala/realpath"
import "io/ioutil"
import "os"
import "path"
import "strings"

func readNPMOffers(packagePath string) ([]Offer, error) {
	var returned []Offer
	nodeModules := path.Join(packagePath, "node_modules")
	directoryEntries, err := readAndStatDir(nodeModules)
	if err != nil {
		if os.IsNotExist(err) {
			return []Offer{}, nil
		}
		return nil, err
	}
	processOffer := func(directory string, scope *string) error {
		anyNewOffers := false
		offers, err := readPackageJSON(directory)
		if err != nil {
			return err
		}
		realDirectory, err := realpath.Realpath(directory)
		if err != nil {
			directory = realDirectory
		}
		for _, offer := range offers {
			if alreadyHaveOffer(returned, offer.OfferID) {
				continue
			}
			anyNewOffers = true
			offer.Artifact.Path = directory
			if scope != nil {
				offer.Artifact.Scope = *scope
			}
			returned = append(returned, offer)
		}
		if anyNewOffers {
			below, recursionError := readNPMOffers(directory)
			if recursionError != nil {
				return recursionError
			}
			returned = append(returned, below...)
		}
		return nil
	}

	for _, entry := range directoryEntries {
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
				err := processOffer(directory, &scope)
				if err != nil {
					return nil, err
				}
			}
		} else { // ./node_modules/package/
			directory := path.Join(nodeModules, name)
			err := processOffer(directory, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	return returned, nil
}

func readPackageJSON(directory string) ([]Offer, error) {
	packageJSON := path.Join(directory, "package.json")
	data, err := ioutil.ReadFile(packageJSON)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Name        string      `json:"name"`
		Version     string      `json:"version"`
		LicenseZero interface{} `json:"licensezero"`
	}
	json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	result, err := interpretLicenseZeroMetadata(parsed.LicenseZero)
	if err != nil {
		return nil, err
	}
	for _, offer := range result {
		offer.Artifact.Type = "npm"
		offer.Artifact.Name = parsed.Name
		offer.Artifact.Version = parsed.Version
	}
	return result, nil
}

func alreadyHaveOffer(offers []Offer, offerID string) bool {
	for _, other := range offers {
		if other.OfferID == offerID {
			return true
		}
	}
	return false
}

func findNPMPackageInfo(directoryPath string) *Offer {
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
	return &Offer{
		Artifact: ArtifactData{
			Type:    "npm",
			Name:    name,
			Scope:   scope,
			Version: parsed.Version,
		},
	}
}
