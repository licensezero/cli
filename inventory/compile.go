package inventory

import (
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
	"os"
	"path"
)

// Inventory describes offers to license artifacts in a working
// directory, grouped into categories reflecting whether the user has or
// needs a license.
type Inventory struct {
	Licensable []Item
	Licensed   []Item
	Own        []Item
	Unlicensed []Item
	Ignored    []Item
	Invalid    []Item
}

// Item describes an offer to license work in an artifact.
type Item struct {
	Type    string
	Path    string
	Scope   string
	Name    string
	Version string
	Public  string
	API     string
	OfferID string
	Offer   api.Offer
}

// Finding represents information about an artifact that references offers.
type Finding struct {
	Type    string
	Path    string
	Scope   string
	Name    string
	Version string
	Public  string
	API     string
	OfferID string
}

func addArtifactOfferToFinding(offer *ArtifactOffer, finding *Finding) {
	finding.API = offer.API
	finding.OfferID = offer.OfferID
	finding.Public = offer.Public
}

// CompileInventory is the top-level function for finding License Zero
// dependencies of a project.
func CompileInventory(
	configPath string,
	cwd string,
	ignoreNoncommercial bool,
	ignoreReciprocal bool,
) (inventory Inventory, err error) {
	// TODO: Don't ignore receipt read errors.
	receipts, _, err := user.ReadReceipts(configPath)
	if err != nil {
		return
	}
	accounts, err := user.ReadAccounts(configPath)
	findings, err := findArtifacts(cwd)
	if err != nil {
		return
	}
	for _, finding := range findings {
		// TODO: Add offer data from APIs to inventory results.
		/*
			offer, err := getOffer(finding.API, finding.OfferID)
			var item Item
			if err != nil {
				inventory.Invalid = append(inventory.Invalid, Item{
					Type:    finding.Type,
					Path:    finding.Path,
					Scope:   finding.Scope,
					Name:    finding.Name,
					Version: finding.Version,
					Public:  finding.Public,
				})
				continue
			} else {
				// See below.
			}
		*/
		item := Item{
			Type:    finding.Type,
			Path:    finding.Path,
			Scope:   finding.Scope,
			Name:    finding.Name,
			Version: finding.Version,
			Public:  finding.Public,
			API:     finding.API,
			OfferID: finding.OfferID,
			// Offer:   *offer,
		}
		inventory.Licensable = append(inventory.Licensable, item)
		if haveReceipt(&item, receipts) {
			inventory.Licensed = append(inventory.Licensed, item)
			continue
		}
		if ownProject(&item, accounts) {
			inventory.Own = append(inventory.Own, item)
			continue
		}
		licenseType := licenseTypeOf(item.Public)
		if (licenseType == noncommercial) && ignoreNoncommercial {
			inventory.Ignored = append(inventory.Ignored, item)
			continue
		}
		if (licenseType == reciprocal) && ignoreReciprocal {
			inventory.Ignored = append(inventory.Ignored, item)
			continue
		}
		inventory.Unlicensed = append(inventory.Unlicensed, item)
	}
	return
}

// findArtifacts calls functions for all the ways the CLI knows to find
// dependencies and other artifacts within a project and combines their
// findings into a single slice.
func findArtifacts(cwd string) (findings []*Finding, err error) {
	finders := []func(string) ([]*Finding, error){
		findLicenseZeroFiles,
		findCargoCrates,
		findRubyGems,
		findGoDeps,
	}
	for _, finder := range finders {
		found, err := finder(cwd)
		if err != nil {
			continue
		}
		for _, finding := range found {
			if alreadyFound(findings, finding) {
				continue
			}
			findings = append(findings, finding)
		}
	}
	return
}

func alreadyFound(findings []*Finding, finding *Finding) bool {
	api := finding.API
	offerID := finding.OfferID
	for _, other := range findings {
		if other.API == api && other.OfferID == offerID {
			return true
		}
	}
	return false
}

func haveReceipt(item *Item, receipts []*user.Receipt) bool {
	api := item.API
	offerID := item.OfferID
	for _, account := range receipts {
		if account.API == api && account.OfferID == offerID {
			return true
		}
	}
	return false
}

func ownProject(item *Item, accounts []*user.Account) bool {
	api := item.API
	licensorID := item.Offer.LicensorID
	for _, account := range accounts {
		if account.API == api && account.LicensorID == licensorID {
			return true
		}
	}
	return false
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

func isSymlink(entry os.FileInfo) bool {
	return entry.Mode()&os.ModeSymlink != 0
}