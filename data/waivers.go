package data

import "encoding/json"
import "io/ioutil"
import "os"
import "path"

type WaiverEnvelope struct {
	Manifest WaiverManifest `json:"manifest"`
}

type WaiverManifest struct {
	ProjectID    string `json:"projectID"`
	Date         string
	Term         string
	Beneficiary  string
	Jurisdiction string
	EMail        string
}

func WaiverPath(home string, projectID string) string {
	return path.Join(home, "waivers", projectID+".json")
}

func ReadWaivers(home string) ([]WaiverManifest, error) {
	directoryPath := path.Join(configPath(home), "waivers")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []WaiverManifest{}, nil
		} else {
			return nil, directoryReadError
		}
	}
	var returned []WaiverManifest
	for _, entry := range entries {
		name := entry.Name()
		waiver, err := readWaiver(home, name)
		if err != nil {
			return nil, err
		}
		if Unexpired(&waiver.Manifest) {
			returned = append(returned, waiver.Manifest)
		}
	}
	return returned, nil
}

func Unexpired(waiver *WaiverManifest) bool {
	// TODO
	return true
}

func readWaiver(home string, file string) (*WaiverEnvelope, error) {
	filePath := path.Join(home, "waivers", file)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var waiver WaiverEnvelope
	json.Unmarshal(data, &waiver)
	return &waiver, nil
}
