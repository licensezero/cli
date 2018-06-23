package data

import "encoding/json"
import "errors"
import "io/ioutil"
import "os"
import "path"
import "strconv"
import "time"

type WaiverEnvelope struct {
	ProjectID string `json:"projectID"`
	Manifest  struct {
		Form        string `json:"FORM"`
		Version     string `json:"VERSION"`
		Beneficiary struct {
			Name         string `json:"name"`
			Jurisdiction string `json:"jurisdiction"`
		} `json:"beneficiary"`
		Date     string `json:"date"`
		Licensor struct {
			Name         string `json:"name"`
			Jurisdiction string `json:"jurisdiction"`
		} `json:"licensor"`
		Project struct {
			ProjectID   string `json:"projectID"`
			Description string `json:"description"`
			Homepage    string `json:"homepage"`
		} `json:"project"`
		Term string `json:"term"`
	} `json:"manifest"`
	Document  string `json:"document"`
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

func WaiversPath(home string) string {
	return path.Join(configPath(home), "waivers")
}

func WaiverPath(home string, projectID string) string {
	return path.Join(WaiversPath(home), projectID+".json")
}

func ReadWaivers(home string) ([]WaiverEnvelope, error) {
	directoryPath := WaiversPath(home)
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
		filePath := path.Join(directoryPath, name)
		waiver, err := ReadWaiver(filePath)
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

func ReadWaiver(filePath string) (*WaiverEnvelope, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var waiver WaiverEnvelope
	json.Unmarshal(data, &waiver)
	return &waiver, nil
}

func WriteWaiver(home string, waiver *WaiverEnvelope) error {
	json, err := json.Marshal(waiver)
	if err != nil {
		return err
	}
	err = os.MkdirAll(WaiversPath(home), 0644)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(WaiverPath(home, waiver.ProjectID), json, 0644)
}
