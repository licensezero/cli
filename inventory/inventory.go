package inventory

import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "os"
import "path"

// Offer describes a License Zero contribution set in inventory.
type Offer struct {
	Type     string                `json:"type"`
	Path     string                `json:"path"`
	Scope    string                `json:"scope,omitempty"`
	Name     string                `json:"name"`
	Version  string                `json:"version"`
	Envelope OfferManifestEnvelope `json:"envelope"`
}

// OfferManifestEnvelope describes a signed offer manifest.
type OfferManifestEnvelope struct {
	LicensorSignature string        `json:"licensorSignature"`
	AgentSignature    string        `json:"agentSignature"`
	Manifest          OfferManifest `json:"license"`
}

// OfferManifest describes contribution set data from licensezero.json.
type OfferManifest struct {
	// Note: These declaration must appear in the order so as to
	// serialize in the correct order for signature verification.
	Repository   string `json:"homepage"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
	OfferID      string `json:"offerID"`
	PublicKey    string `json:"publicKey"`
	Terms        string `json:"terms"`
	Version      string `json:"version"`
}

// Offers describes the categorization of offers in inventory.
type Offers struct {
	Licensable []Offer `json:"licensable"`
	Licensed   []Offer `json:"licensed"`
	Waived     []Offer `json:"waived"`
	Unlicensed []Offer `json:"unlicensed"`
	Ignored    []Offer `json:"ignored"`
	Invalid    []Offer `json:"invalid"`
}

// Inventory finds License Zero offers included or referenced in a working tree.
func Inventory(home string, cwd string, ignoreNC bool, ignoreR bool) (*Offers, error) {
	identity, _ := data.ReadIdentity(home)
	var licenses []data.LicenseEnvelope
	var waivers []data.WaiverEnvelope
	if identity != nil {
		readLicenses, err := data.ReadLicenses(home)
		if err != nil {
			licenses = readLicenses
		}
		readWaivers, err := data.ReadWaivers(home)
		if err != nil {
			waivers = readWaivers
		}
	}
	offers, err := readOffers(cwd)
	if err != nil {
		return nil, err
	}
	agentPublicKey, err := api.FetchAgentPublicKey()
	if err != nil {
		return nil, err
	}
	var returned Offers
	for _, result := range offers {
		licensor, err := api.Read(result.Envelope.Manifest.OfferID)
		if err != nil {
			returned.Invalid = append(returned.Invalid, result)
			continue
		}
		err = CheckMetadata(&result.Envelope, licensor.PublicKey, agentPublicKey)
		if err != nil {
			returned.Invalid = append(returned.Invalid, result)
			continue
		} else {
			returned.Licensable = append(returned.Licensable, result)
		}
		if haveLicense(&result, licenses, identity) {
			returned.Licensed = append(returned.Licensed, result)
			continue
		}
		if haveWaiver(&result, waivers, identity) {
			returned.Waived = append(returned.Waived, result)
			continue
		}
		if identity != nil && ownOffer(&result, identity) {
			continue
		}
		terms := result.Envelope.Manifest.Terms
		if (terms == "noncommercial" || terms == "prosperity") && ignoreNC {
			returned.Ignored = append(returned.Ignored, result)
			continue
		}
		if (terms == "reciprocal" || terms == "parity") && ignoreR {
			returned.Ignored = append(returned.Ignored, result)
			continue
		}
		returned.Unlicensed = append(returned.Unlicensed, result)
	}
	return &returned, nil
}

func haveLicense(offer *Offer, licenses []data.LicenseEnvelope, identity *data.Identity) bool {
	for _, license := range licenses {
		if license.OfferID == offer.Envelope.Manifest.OfferID &&
			license.Manifest.Licensee.Name == identity.Name &&
			license.Manifest.Licensee.Jurisdiction == identity.Jurisdiction &&
			license.Manifest.Licensee.EMail == identity.EMail {
			return true
		}
	}
	return false
}

func haveWaiver(offer *Offer, waivers []data.WaiverEnvelope, identity *data.Identity) bool {
	for _, waiver := range waivers {
		if waiver.OfferID == offer.Envelope.Manifest.OfferID &&
			waiver.Manifest.Beneficiary.Name == identity.Name &&
			waiver.Manifest.Beneficiary.Jurisdiction == identity.Jurisdiction {
			return true
		}
	}
	return false
}

func ownOffer(offer *Offer, identity *data.Identity) bool {
	return offer.Envelope.Manifest.Name == identity.Name &&
		offer.Envelope.Manifest.Jurisdiction == identity.Jurisdiction
}

func readOffers(cwd string) ([]Offer, error) {
	descenders := []func(string) ([]Offer, error){
		readNPMOffers,
		readRubyGemsOffers,
		readGoDeps,
		readCargoCrates,
		recurseLicenseZeroFiles,
	}
	returned := []Offer{}
	for _, descender := range descenders {
		offers, err := descender(cwd)
		if err == nil {
			for _, offer := range offers {
				offerID := offer.Envelope.Manifest.OfferID
				if alreadyHaveOffer(returned, offerID) {
					continue
				}
				returned = append(returned, offer)
			}
		}
	}
	return returned, nil
}

func isSymlink(entry os.FileInfo) bool {
	return entry.Mode()&os.ModeSymlink != 0
}

// Like ioutil.ReadDir, but don't sort, and read all symlinks.
func readAndStatDir(directoryPath string) ([]os.FileInfo, error) {
	directory, err := os.Open(directoryPath)
	if err != nil {
		return nil, err
	}
	entries, err := directory.Readdir(-1)
	directory.Close()
	if err != nil {
		return nil, err
	}
	returned := make([]os.FileInfo, len(entries))
	for i, entry := range entries {
		if isSymlink(entry) {
			linkPath := path.Join(directoryPath, entry.Name())
			targetPath, err := os.Readlink(linkPath)
			if err != nil {
				return nil, err
			}
			if !path.IsAbs(targetPath) {
				targetPath = path.Join(path.Dir(directoryPath), targetPath)
			}
			newEntry, err := os.Stat(targetPath)
			if err != nil {
				return nil, err
			}
			returned[i] = newEntry
		} else {
			returned[i] = entry
		}
	}
	return returned, nil
}
