package inventory

import "fmt"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "os"
import "path"

// Offer describes a License Zero offer in inventory.
type Offer struct {
	Artifact   ArtifactData `json:"artifact"`
	License    LicenseData  `json:"license"`
	OfferID    string       `json:"offerID"`
	LicensorID string       `json:"licensorID"`
	Envelope   Verifiable   `json:"omit"`
}

// ArtifactData groups data about package, crate, and so on.
type ArtifactData struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Scope   string `json:"scope,omitempty"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// LicenseData groups data about a public license for an artifact.
type LicenseData struct {
	Terms   string `json:"terms"`
	Version string `json:"version"`
}

// Envelope describes a signed data envelope with offer information.
type Envelope interface {
	offer() Offer
}

// Offers describes the categorization of projects in inventory.
type Offers struct {
	Licensable []Offer `json:"licensable"`
	Licensed   []Offer `json:"licensed"`
	Waived     []Offer `json:"waived"`
	Unlicensed []Offer `json:"unlicensed"`
	Ignored    []Offer `json:"ignored"`
	Invalid    []Offer `json:"invalid"`
}

// Inventory finds License Zero projects included or referenced in a working tree.
func Inventory(home string, cwd string, ignoreNC bool, ignoreR bool) (*Offers, error) {
	identity, _ := data.ReadIdentity(home)
	var licenses []data.License
	var waivers []data.Waiver
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
	licensor, _ := data.ReadLicensor(home)
	offers, err := readOffers(cwd)
	if err != nil {
		return nil, err
	}
	agentPublicKey, err := api.FetchAgentPublicKey()
	if err != nil {
		return nil, err
	}
	var returned Offers
	for _, offer := range offers {
		offerLicensor, err := api.Read(offer.OfferID)
		if err != nil {
			returned.Invalid = append(returned.Invalid, offer)
			continue
		}
		err = offer.Envelope.verifyLicensorSignature(offerLicensor.PublicKey)
		if err != nil {
			fmt.Printf("invalid licensor signature")
			returned.Invalid = append(returned.Invalid, offer)
			continue
		}
		err = offer.Envelope.verifyAgentSignature(agentPublicKey)
		if err != nil {
			returned.Invalid = append(returned.Invalid, offer)
			continue
		}
		returned.Licensable = append(returned.Licensable, offer)
		if haveLicense(&offer, licenses, identity) {
			returned.Licensed = append(returned.Licensed, offer)
			continue
		}
		if haveWaiver(&offer, waivers, identity) {
			returned.Waived = append(returned.Waived, offer)
			continue
		}
		if licensor != nil && ownOffer(&offer, licensor) {
			continue
		}
		terms := offer.License.Terms
		if (terms == "noncommercial" || terms == "prosperity") && ignoreNC {
			returned.Ignored = append(returned.Ignored, offer)
			continue
		}
		if (terms == "reciprocal" || terms == "parity") && ignoreR {
			returned.Ignored = append(returned.Ignored, offer)
			continue
		}
		returned.Unlicensed = append(returned.Unlicensed, offer)
	}
	return &returned, nil
}

func haveLicense(offer *Offer, licenses []data.License, identity *data.Identity) bool {
	for _, license := range licenses {
		if license.OfferID == offer.OfferID &&
			license.LicenseeName == identity.Name &&
			license.LicenseeJurisdiction == identity.Jurisdiction &&
			license.LicenseeEMail == identity.EMail {
			return true
		}
	}
	return false
}

func haveWaiver(offer *Offer, waivers []data.Waiver, identity *data.Identity) bool {
	for _, waiver := range waivers {
		// TODO: Also compare beneificary e-mails.
		if waiver.OfferID == offer.OfferID &&
			waiver.BeneficiaryName == identity.Name &&
			waiver.BeneficiaryJurisdiction == identity.Jurisdiction {
			return true
		}
	}
	return false
}

func ownOffer(offer *Offer, licensor *data.Licensor) bool {
	return offer.LicensorID == licensor.LicensorID
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
				if alreadyHaveOffer(returned, offer.OfferID) {
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
