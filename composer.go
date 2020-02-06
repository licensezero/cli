package cli

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

func readComposerPackageMetadata(directoryPath string) *Finding {
	composerJSON := path.Join(directoryPath, "composer.json")
	data, err := ioutil.ReadFile(composerJSON)
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
	return &Finding{
		Type:    "composer",
		Name:    parsed.Name,
		Version: parsed.Version,
	}
}
