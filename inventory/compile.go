package inventory

import (
	"net/http"
	"os"
	"path"

	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
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
	Server  string
	OfferID string
	Offer   api.Offer
	Seller  api.Seller
}

// Finding represents information about an artifact that references offers.
type Finding struct {
	Type    string
	Path    string
	Scope   string
	Name    string
	Version string
	Public  string
	Server  string
	OfferID string
}

func addArtifactOfferToFinding(offer *ArtifactOffer, finding *Finding) {
	finding.Server = offer.Server
	finding.OfferID = offer.OfferID
	finding.Public = offer.Public
}

// Compile is the top-level function for finding License Zero
// dependencies of a project.
func Compile(
	configPath string,
	cwd string,
	ignoreNoncommercial bool,
	ignoreReciprocal bool,
	client *http.Client,
) (inventory Inventory, err error) {
	receipts, _, err := user.ReadReceipts()
	if err != nil {
		return
	}
	accounts, err := user.ReadAccounts()
	if err != nil {
		return
	}
	findings, err := findArtifacts(cwd)
	if err != nil {
		return
	}
	for _, finding := range findings {
		brokerServer := api.BrokerServer{
			Client: client,
			Base:   finding.Server,
		}
		offer, err := brokerServer.Offer(finding.OfferID)
		var item Item
		handleInvalid := func() {
			inventory.Invalid = append(inventory.Invalid, Item{
				Type:    finding.Type,
				Path:    finding.Path,
				Scope:   finding.Scope,
				Name:    finding.Name,
				Version: finding.Version,
				Public:  finding.Public,
			})
		}
		if err != nil {
			handleInvalid()
			continue
		}
		seller, err := brokerServer.Seller(offer.SellerID)
		if err != nil {
			handleInvalid()
			continue
		}
		item = Item{
			Type:    finding.Type,
			Path:    finding.Path,
			Scope:   finding.Scope,
			Name:    finding.Name,
			Version: finding.Version,
			Public:  finding.Public,
			Server:  finding.Server,
			OfferID: finding.OfferID,
			Offer:   *offer,
			Seller:  *seller,
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
		if (item.Public == "noncommercial") && ignoreNoncommercial {
			inventory.Ignored = append(inventory.Ignored, item)
			continue
		}
		if (item.Public == "share alike") && ignoreReciprocal {
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
	server := finding.Server
	offerID := finding.OfferID
	for _, other := range findings {
		if other.Server == server && other.OfferID == offerID {
			return true
		}
	}
	return false
}

func haveReceipt(item *Item, receipts []*api.Receipt) bool {
	server := item.Server
	offerID := item.OfferID
	for _, account := range receipts {
		if account.License.Values.Server == server &&
			account.License.Values.OfferID == offerID {
			return true
		}
	}
	return false
}

func ownProject(item *Item, accounts []*user.Account) bool {
	server := item.Server
	sellerID := item.Offer.SellerID
	for _, account := range accounts {
		if account.Server == server && account.SellerID == sellerID {
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
