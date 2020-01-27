package inventory

import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "os"
import "path"

// ArtifactMetadata describes offer metadata in an artifact.
type ArtifactMetadata struct {
	Schema string           `json:"schema"`
	Offers []OfferReference `json:"offers"`
}

// DescenderResult describes an OfferReference read from a project.
type DescenderResult struct {
	Type    string
	Path    string
	Name    string
	Scope   string
	Version string
	Offer   OfferReference
}

// OfferReference describes a reference to a license offer for an artifact.
type OfferReference struct {
	OfferID       string `json:"id"`
	API           string `json:"api"`
	PublicLicense string `json:"public"`
}

// Inventory categorizes license offers for artifacts ina project.
type Inventory struct {
	Licensable []Result
	Licensed   []Result
	Unlicensed []Result
	Ignored    []Result
	Invalid    []Result
}

// Result combines data about the reference read from the project and the offer data received from vendor APIs.
type Result struct {
	Local  DescenderResult
	Remote api.Offer
}

// ReadInventory finds artifacts and categorizes artifacts included or referenced in a project.
func ReadInventory(home string, cwd string, ignoreNC bool, ignoreR bool) (*Inventory, error) {
	identity, _ := data.ReadIdentity(home)
	var receipts []data.Receipt
	if identity != nil {
		readLicenses, err := data.ReadReceipts(home)
		if err != nil {
			receipts = readLicenses
		}
	}
	licensorAccounts, _ := data.ReadLicensorAccounts(home)
	descenderResults, err := descend(cwd)
	if err != nil {
		return nil, err
	}
	var returned Inventory
	for _, descenderResult := range descenderResults {
		offer, err := api.GetOffer(descenderResult.Offer.API, descenderResult.Offer.OfferID)
		if err != nil {
			returned.Invalid = append(returned.Invalid, Result{
				Local: descenderResult,
			})
			continue
		}
		result := Result{
			Local:  descenderResult,
			Remote: *offer,
		}
		if haveReceipt(&descenderResult.Offer, receipts) {
			returned.Licensed = append(returned.Licensed, result)
			continue
		}
		if identity != nil && ownOffer(offer, licensorAccounts) {
			continue
		}
		license := descenderResult.Offer.PublicLicense
		licenseType := TypeOfLicense(license)
		if (licenseType == Noncommercial) && ignoreNC {
			returned.Ignored = append(returned.Ignored, result)
			continue
		}
		if (licenseType == Reciprocal) && ignoreR {
			returned.Ignored = append(returned.Ignored, result)
			continue
		}
		returned.Unlicensed = append(returned.Unlicensed, result)
	}
	return &returned, nil
}

func haveReceipt(offerReference *OfferReference, receipts []data.Receipt) bool {
	for _, receipt := range receipts {
		if receipt.License.Values.Vendor.API == offerReference.API &&
			receipt.License.Values.Offer == offerReference.OfferID {
			return true
		}
	}
	return false
}

func ownOffer(offer *api.Offer, accounts []data.LicensorAccount) bool {
	for _, account := range accounts {
		if account.API == offer.API &&
			account.LicensorID == offer.LicensorID {
			return true
		}
	}
	return false
}

func descend(cwd string) ([]DescenderResult, error) {
	descenders := []func(string) ([]DescenderResult, error){
		readNPMProjects,
		readRubyGemsProjects,
		readGoDeps,
		readCargoCrates,
		recurseLicenseZeroFiles,
	}
	returned := []DescenderResult{}
	for _, descender := range descenders {
		results, err := descender(cwd)
		if err == nil {
			for _, result := range results {
				if alreadyHave(returned, &result.Offer) {
					continue
				}
				returned = append(returned, result)
			}
		}
	}
	return returned, nil
}

func alreadyHave(have []DescenderResult, found *OfferReference) bool {
	for _, other := range have {
		if other.Offer.API == found.API &&
			other.Offer.OfferID == found.OfferID {
			return true
		}
	}
	return false
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
