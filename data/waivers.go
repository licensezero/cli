package data

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "os"
import "path"
import "strconv"
import "time"

// WaiverEnvelope describes a fully parsed waiver file.
type WaiverEnvelope struct {
	ProjectID      string
	Manifest       WaiverManifest
	ManifestString string
	Document       string
	Signature      string
	PublicKey      string
}

// WaiverFile describes a partially parsed waiver file.
type WaiverFile struct {
	ProjectID string `json:"projectID"`
	Manifest  string `json:"manifest"`
	Document  string `json:"document"`
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

// WaiverManifest describes signed waiver data.
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

func waiversPath(home string) string {
	return path.Join(ConfigPath(home), "waivers")
}

func waiverPath(home string, projectID string) string {
	return path.Join(waiversPath(home), projectID+".json")
}

// ReadWaivers reads all waivers from the configuration directory.
func ReadWaivers(home string) ([]WaiverEnvelope, error) {
	directoryPath := waiversPath(home)
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []WaiverEnvelope{}, nil
		}
		return nil, directoryReadError
	}
	var returned []WaiverEnvelope
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(directoryPath, name)
		waiver, err := ReadWaiver(filePath)
		if err != nil {
			return nil, err
		}
		unexpired, _ := checkUnexpired(waiver)
		if unexpired {
			returned = append(returned, *waiver)
		}
	}
	return returned, nil
}

func checkUnexpired(waiver *WaiverEnvelope) (bool, error) {
	termString := waiver.Manifest.Term
	if termString == "forever" {
		return true, nil
	}
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
	return true, nil
}

// ReadWaiver reads a waiver file.
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

// WriteWaiver writes a waiver to the configuration directory
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
	err = os.MkdirAll(waiversPath(home), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(waiverPath(home, waiver.ProjectID), json, 0644)
}

// CheckWaiverSignature verifies the signatures to a waiver.
func CheckWaiverSignature(waiver *WaiverEnvelope, publicKey string) error {
	serialized, err := json.Marshal(waiver.Manifest)
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return errors.New("could not serialize manifest")
	}
	if waiver.ProjectID != waiver.Manifest.Project.ProjectID {
		return errors.New("project IDs do not match")
	}
	err = checkSignature(
		publicKey,
		waiver.Signature,
		[]byte(waiver.ManifestString+"\n\n"+waiver.Document),
	)
	if err != nil {
		return err
	}
	return nil
}

func compactWaiverManifest(data *WaiverManifest) (*bytes.Buffer, error) {
	serialized, err := json.Marshal(data)
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return nil, errors.New("could not serialize")
	}
	return compacted, nil
}
