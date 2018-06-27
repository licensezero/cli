package data

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "os"
import "path"
import "strconv"
import "time"

type WaiverEnvelope struct {
	ProjectID      string
	Manifest       WaiverManifest
	ManifestString string
	Document       string
	Signature      string
	PublicKey      string
}

type WaiverFile struct {
	ProjectID string `json:"projectID"`
	Manifest  string `json:"manifest"`
	Document  string `json:"document"`
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

type WaiverManifest struct {
	Form        string `json:"FORM"`
	Version     string `json:"VERSION"`
	Beneficiary struct {
		Jurisdiction string `json:"jurisdiction"`
		Name         string `json:"name"`
	} `json:"beneficiary"`
	Date     string `json:"date"`
	Licensor struct {
		Jurisdiction string `json:"jurisdiction"`
		Name         string `json:"name"`
	} `json:"licensor"`
	Project struct {
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
		ProjectID   string `json:"projectID"`
	} `json:"project"`
	Term string `json:"term"`
}

func WaiversPath(home string) string {
	return path.Join(ConfigPath(home), "waivers")
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
	var file WaiverFile
	json.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}
	var manifest WaiverManifest
	err = json.Unmarshal([]byte(file.Manifest), &manifest)
	if err != nil {
		return nil, err
	}
	return &WaiverEnvelope{
		Manifest:       manifest,
		ManifestString: file.Manifest,
		ProjectID:      file.ProjectID,
		Document:       file.Document,
		PublicKey:      file.PublicKey,
		Signature:      file.Signature,
	}, nil
}

func WriteWaiver(home string, waiver *WaiverEnvelope) error {
	file := WaiverFile{
		Manifest:  waiver.ManifestString,
		ProjectID: waiver.ProjectID,
		Document:  waiver.Document,
		PublicKey: waiver.PublicKey,
		Signature: waiver.Signature,
	}
	json, err := json.Marshal(file)
	if err != nil {
		return err
	}
	err = os.MkdirAll(WaiversPath(home), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(WaiverPath(home, waiver.ProjectID), json, 0644)
}

func CheckWaiverSignature(waiver *WaiverEnvelope, publicKey string) error {
	serialized, err := json.Marshal(waiver.Manifest)
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return errors.New("Could not serialize manifest.")
	}
	if waiver.ProjectID != waiver.Manifest.Project.ProjectID {
		return errors.New("Project IDs do not match.")
	}
	err = checkSignature(
		publicKey,
		waiver.Signature,
		[]byte(waiver.ManifestString+"\n\n"+waiver.Document),
	)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func compactWaiverManifest(data *WaiverManifest) (*bytes.Buffer, error) {
	serialized, err := json.Marshal(data)
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return nil, errors.New("Could not serialize.")
	}
	return compacted, nil
}
