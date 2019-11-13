package inventory

import "bytes"
import "os/exec"
import "regexp"
import "strings"

func readRubyGemsOffers(packagePath string) ([]Offer, error) {
	// Run `bundle show` in the working directory to list dependencies.
	showNamesAndVersions := exec.Command("bundle", "show")
	showNamesAndVersions.Dir = packagePath
	var first bytes.Buffer
	showNamesAndVersions.Stdout = &first
	err := showNamesAndVersions.Run()
	if err != nil {
		return nil, err
	}
	namesAndVersions := strings.Split(string(first.Bytes()), "\n")
	// Run `bundle show --paths` to list dependencies' paths.
	showPaths := exec.Command("bundle", "show", "--paths")
	var second bytes.Buffer
	showPaths.Stdout = &second
	err = showPaths.Run()
	if err != nil {
		return nil, err
	}
	paths := strings.Split(string(second.Bytes()), "\n")
	var returned []Offer
	// Parse each line of output.
	re, _ := regexp.Compile(`^\s+\*\s+([^(]+) \((.+)\)$`)
	for i, line := range namesAndVersions[1:] {
		result := re.FindStringSubmatch(line)
		if len(result) == 0 {
			continue
		}
		name := result[1]
		version := result[2]
		gemPath := paths[i]
		// Try to read a licensezero.json file there.
		offers, err := ReadLicenseZeroJSON(gemPath)
		if err != nil {
			continue
		}
		for _, offer := range offers {
			if alreadyHaveOffer(returned, offer.Envelope.Manifest.OfferID) {
				continue
			}
			offer.Type = "rubygem"
			offer.Name = name
			offer.Version = version
			returned = append(returned, offer)
		}
	}
	return returned, nil
}
