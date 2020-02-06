package cli

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
)

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
