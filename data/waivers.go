package data

import "encoding/json"
import "errors"
import "io/ioutil"
import "os"
import "path"
import "strconv"
import "time"

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

func ReadWaivers(home string) ([]WaiverEnvelope, error) {
	directoryPath := path.Join(configPath(home), "waivers")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []WaiverEnvelope{}, nil
		} else {
			return nil, directoryReadError
		}
	}
	var returned []WaiverEnvelope
	for _, entry := range entries {
		name := entry.Name()
		waiver, err := readWaiver(home, name)
		if err != nil {
			return nil, err
		}
		unexpired, _ := Unexpired(waiver)
		if unexpired {
			returned = append(returned, *waiver)
		}
	}
	return returned, nil
}

func Unexpired(waiver *WaiverEnvelope) (bool, error) {
	termString := waiver.Manifest.Term
	if termString == "forever" {
		return true, nil
	} else {
		days, err := strconv.Atoi(termString)
		if err != nil {
			return false, errors.New("could not parse term")
		}
		expiration, err := time.Parse(time.RFC3339, waiver.Manifest.Date)
		if err != nil {
			return false, err
		}
		expiration.AddDate(0, 0, days)
		return expiration.After(time.Now()), nil
	}
	return true, nil
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
