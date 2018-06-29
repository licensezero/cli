package subcommands

import "encoding/json"
import "errors"
import "flag"
import "github.com/licensezero/cli/inventory"
import "io/ioutil"
import "os"
import "strings"

const readmeDescription = "Append licensing information to README."

var README = Subcommand{
	Tag:         "seller",
	Description: readmeDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("readme", flag.ExitOnError)
		silent := Silent(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = readmeUsage
		flagSet.Parse(args)
		var existing string
		checkForLegacyPackageJSON(paths.CWD)
		data, err := ioutil.ReadFile("README.md")
		if err != nil {
			if os.IsNotExist(err) {
				existing = ""
			} else {
				Fail("Error reading README.md.")
			}
		} else {
			existing = string(data)
		}
		projectIDs, termsIDs, err := readEntries(paths.CWD)
		if err != nil {
			Fail(err.Error())
		}
		if len(existing) > 0 {
			existing = existing + "\n\n"
		}
		existing = existing + "# Licensing"
		if len(projectIDs) == 0 {
			Fail("No License Zero project metadata in licensezero.json.")
		}
		haveReciprocal := false
		haveNoncommercial := false
		haveParity := false
		haveProsperity := false
		for _, terms := range termsIDs {
			if terms == "noncommercial" {
				haveNoncommercial = true
			} else if terms == "reciprocal" {
				haveReciprocal = true
			} else if terms == "parity" {
				haveParity = true
			} else if terms == "prosperity" {
				haveProsperity = true
			}
		}
		multiple := twoOrMore([]bool{haveReciprocal, haveNoncommercial, haveParity, haveProsperity})
		var licenseScope string
		if multiple {
			licenseScope = "Some contributions to this package "
		} else {
			licenseScope = "This package "
		}
		summaries := []string{}
		availabilities := []string{}
		if haveReciprocal {
			summaries = append(
				summaries,
				licenseScope+
					"is free to use in open source under the terms of "+
					"the [License Zero Reciprocal Public License](./LICENSE).",
			)
			availabilities = append(
				availabilities,
				"Licenses for use in closed software "+
					"are available via [licensezero.com](https://licensezero.com).",
			)
		} else if haveNoncommercial {
			summaries = append(
				summaries,
				licenseScope+
					"is free to use for commercial purposes for a trial period under the terms of "+
					"the [License Zero Noncommercial Public License](./LICENSE).",
			)
			availabilities = append(
				availabilities,
				"Licenses for long-term commercial use "+
					"are available via [licensezero.com](https://licensezero.com).",
			)
		} else if haveParity {
			summaries = append(
				summaries,
				licenseScope+
					"is free to use in open source under the terms of "+
					"[Parity Public License](./LICENSE).",
			)
			availabilities = append(
				availabilities,
				"Licenses for use in closed software "+
					"are available via [licensezero.com](https://licensezero.com).",
			)
		} else if haveProsperity {
			summaries = append(
				summaries,
				licenseScope+
					"is free to use for commercial purposes for a trial period under the terms of "+
					"[The Prosperity Public License](./LICENSE).",
			)
			availabilities = append(
				availabilities,
				"Licenses for long-term commercial use "+
					"are available via [licensezero.com](https://licensezero.com).",
			)
		} else {
			Fail("Unrecognized License Zero project terms.")
		}
		existing = existing + "\n\n" + strings.Join(summaries, "\n\n")
		existing = existing + "\n\n" + strings.Join(availabilities, "\n\n")
		for _, projectID := range projectIDs {
			projectLink := "https://licensezero.com/projects/" + projectID
			badge := "" +
				"[" +
				"![licensezero.com pricing](" + projectLink + "/badge.svg)" +
				"]" +
				"[" + projectLink + "]"
			existing = existing + "\n\n" + badge
		}
		err = ioutil.WriteFile("README.md", []byte(existing), 0644)
		if err != nil {
			Fail("Error writing README.md.")
		}
		if !*silent {
			os.Stdout.WriteString("Wrote to README.md\n")
		}
		os.Exit(0)
	},
}

func twoOrMore(values []bool) bool {
	counter := 0
	for _, value := range values {
		if value {
			counter++
		}
		if counter == 2 {
			return true
		}
	}
	return false
}

func readEntries(directory string) ([]string, []string, error) {
	data, err := ioutil.ReadFile("licensezero.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.New("Could not read licensezero.json.")
		} else {
			return nil, nil, err
		}
	}
	var existingMetadata inventory.LicenseZeroJSONFile
	err = json.Unmarshal(data, &existingMetadata)
	if err != nil {
		return nil, nil, errors.New("Could not parse licensezero.json metadata.")
	}
	// TODO: Validate licensezero.json metadata entries.
	var projectIDs []string
	var terms []string
	for _, entry := range existingMetadata.Envelopes {
		projectIDs = append(projectIDs, entry.Manifest.ProjectID)
		terms = append(terms, entry.Manifest.Terms)
	}
	return projectIDs, terms, nil
}

func readmeUsage() {
	usage := readmeDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero readme\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"silent": silentLine,
		})
	Fail(usage)
}
