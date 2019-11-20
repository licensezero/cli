package data

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "os"
import "path"
import "strconv"
import "time"

// Waiver describes a waiver.
type Waiver struct {
	OfferID                 string
	Date                    string
	Term                    string
	BeneficiaryName         string
	BeneficiaryJurisdiction string
	BeneficiaryEMail        string
}

// WaiverVersion abstracts over waiver versions.
type WaiverVersion interface {
	waiver() Waiver
}

// Version1WaiverEnvelope describes a fully parsed waiver file.
type Version1WaiverEnvelope struct {
	OfferID        string
	Manifest       Version1WaiverManifest
	ManifestString string
	Document       string
	Signature      string
	PublicKey      string
}

func (waiver Version1WaiverEnvelope) waiver() Waiver {
	beneficiary := waiver.Manifest.Beneficiary
	return Waiver{
		OfferID:                 waiver.OfferID,
		BeneficiaryName:         beneficiary.Name,
		BeneficiaryJurisdiction: beneficiary.Jurisdiction,
	}
}

// Version1WaiverFile describes a partially parsed waiver file.
type Version1WaiverFile struct {
	OfferID   string `json:"offerID"`
	Manifest  string `json:"manifest"`
	Document  string `json:"document"`
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

// Version1WaiverManifest describes signed waiver data.
type Version1WaiverManifest struct {
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
	Offer struct {
		Description string `json:"description"`
		Repository  string `json:"homepage"`
		OfferID     string `json:"offerID"`
	} `json:"project"`
	Term string `json:"term"`
}

func waiversPath(home string) string {
	return path.Join(ConfigPath(home), "waivers")
}

func waiverPath(home string, offerID string) string {
	return path.Join(waiversPath(home), offerID+".json")
}

// ReadWaivers reads all waivers from the configuration directory.
func ReadWaivers(home string) ([]Waiver, error) {
	directoryPath := waiversPath(home)
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []Waiver{}, nil
		}
		return nil, directoryReadError
	}
	var returned []Waiver
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(directoryPath, name)
		raw, err := ReadWaiver(filePath)
		waiver := raw.waiver()
		if err != nil {
			return nil, err
		}
		unexpired, _ := checkUnexpired(&waiver)
		if unexpired {
			returned = append(returned, waiver)
		}
	}
	return returned, nil
}

func checkUnexpired(waiver *Waiver) (bool, error) {
	termString := waiver.Term
	if termString == "forever" {
		return true, nil
	}
	days, err := strconv.Atoi(termString)
	if err != nil {
		return false, errors.New("could not parse term")
	}
	expiration, err := time.Parse(time.RFC3339, waiver.Date)
	if err != nil {
		return false, err
	}
	expiration.AddDate(0, 0, days)
	return expiration.After(time.Now()), nil
}

// ReadWaiver reads a waiver file.
func ReadWaiver(filePath string) (*Version1WaiverEnvelope, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var file Version1WaiverFile
	json.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}
	var manifest Version1WaiverManifest
	err = json.Unmarshal([]byte(file.Manifest), &manifest)
	if err != nil {
		return nil, err
	}
	return &Version1WaiverEnvelope{
		Manifest:       manifest,
		ManifestString: file.Manifest,
		OfferID:        file.OfferID,
		Document:       file.Document,
		PublicKey:      file.PublicKey,
		Signature:      file.Signature,
	}, nil
}

// WriteWaiver writes a waiver to the configuration directory
func WriteWaiver(home string, waiver *Version1WaiverEnvelope) error {
	file := Version1WaiverFile{
		Manifest:  waiver.ManifestString,
		OfferID:   waiver.OfferID,
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
	return ioutil.WriteFile(waiverPath(home, waiver.OfferID), json, 0644)
}

// CheckWaiverSignature verifies the signatures to a waiver.
func CheckWaiverSignature(waiver *Version1WaiverEnvelope, publicKey string) error {
	serialized, err := json.Marshal(waiver.Manifest)
	if err != nil {
		return errors.New("coiuld not serialize waiver manifest")
	}
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return errors.New("could not compact waiver manifest")
	}
	if waiver.OfferID != waiver.Manifest.Offer.OfferID {
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

func compactVersion1WaiverManifest(data *Version1WaiverManifest) (*bytes.Buffer, error) {
	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("could not serialize waiver manifest")
	}
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return nil, errors.New("could not compact waiver manifest")
	}
	return compacted, nil
}
