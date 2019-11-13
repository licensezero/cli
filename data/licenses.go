package data

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "os"
import "path"

// LicenseEnvelope describes fully parsed licensezero.json metadata about an offer.
type LicenseEnvelope struct {
	Manifest       LicenseManifest
	ManifestString string
	OfferID        string
	Document       string
	PublicKey      string
	Signature      string
}

// LicenseFile describes partially parsed licensezero.json metadata about a contribution set.
type LicenseFile struct {
	Manifest  string `json:"manifest"`
	OfferID   string `json:"offerID"`
	Document  string `json:"document"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
}

// LicenseManifest describes signed licensezero.json metadata about a contribution set.
type LicenseManifest struct {
	Date     string `json:"date"`
	Form     string `json:"FORM"`
	Licensee struct {
		Name         string `json:"name"`
		Jurisdiction string `json:"jurisdiction"`
		EMail        string `json:"email"`
	}
	Licensor struct {
		Jurisdiction string `json:"jurisdiction"`
		Name         string `json:"name"`
	}
	OrderID string `json:"orderID"`
	Price   int    `json:"price"`
	Offer   struct {
		Description string `json:"description"`
		Repository  string `json:"homepage"`
		OfferID     string `json:"offerID"`
	}
	Version string `json:"VERSION"`
}

func licensePath(home string, offerID string) string {
	return path.Join(licensesPath(home), offerID+".json")
}

func licensesPath(home string) string {
	return path.Join(ConfigPath(home), "licenses")
}

// ReadLicenses reads all saved licenses from the CLI configuration directory.
func ReadLicenses(home string) ([]LicenseEnvelope, error) {
	directoryPath := path.Join(ConfigPath(home), "licenses")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []LicenseEnvelope{}, nil
		}
		return nil, directoryReadError
	}
	var returned []LicenseEnvelope
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(home, "licenses", name)
		license, err := ReadLicense(filePath)
		if err != nil {
			return nil, err
		}
		returned = append(returned, *license)
	}
	return returned, nil
}

// LicenseFileToEnvelope fully parses license file data.
func LicenseFileToEnvelope(file *LicenseFile) (*LicenseEnvelope, error) {
	var manifest LicenseManifest
	err := json.Unmarshal([]byte(file.Manifest), &manifest)
	if err != nil {
		return nil, err
	}
	return &LicenseEnvelope{
		Manifest:       manifest,
		ManifestString: file.Manifest,
		OfferID:        file.OfferID,
		Document:       file.Document,
		PublicKey:      file.PublicKey,
		Signature:      file.Signature,
	}, nil
}

// ReadLicense reads a license file from disk.
func ReadLicense(filePath string) (*LicenseEnvelope, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var file LicenseFile
	err = json.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}
	return LicenseFileToEnvelope(&file)
}

// WriteLicense writes a license file to the CLI configuration directory.
func WriteLicense(home string, license *LicenseEnvelope) error {
	file := LicenseFile{
		Manifest:  license.ManifestString,
		OfferID:   license.OfferID,
		Document:  license.Document,
		PublicKey: license.PublicKey,
		Signature: license.Signature,
	}
	json, err := json.Marshal(file)
	if err != nil {
		return err
	}
	err = os.MkdirAll(licensesPath(home), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(licensePath(home, license.OfferID), json, 0644)
}

// CheckLicenseSignature verifies the signatures to a liecnse envelope.
func CheckLicenseSignature(license *LicenseEnvelope, publicKey string) error {
	serialized, err := json.Marshal(license.Manifest)
	if err != nil {
		return errors.New("could not serialize license manifest")
	}
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return errors.New("could not compact license manifest")
	}
	if license.OfferID != license.Manifest.Offer.OfferID {
		return errors.New("offer IDs do not match")
	}
	err = checkSignature(
		publicKey,
		license.Signature,
		[]byte(license.ManifestString+"\n\n"+license.Document),
	)
	if err != nil {
		return err
	}
	return nil
}

func compactLicenseManifest(data *LicenseManifest) (*bytes.Buffer, error) {
	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("could not serialize license manifest")
	}
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return nil, errors.New("could not compact license manifest")
	}
	return compacted, nil
}
