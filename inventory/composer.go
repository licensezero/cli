package inventory

import "encoding/json"
import "io/ioutil"
import "path"

func findComposerPackageInfo(directoryPath string) *Offer {
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
	return &Offer{
		Artifact: ArtifactData{
			Type:    "composer",
			Name:    parsed.Name,
			Version: parsed.Version,
		},
	}
}
