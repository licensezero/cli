package subcommands

import "encoding/json"
import "errors"
import "flag"
import "io/ioutil"
import "os"

const readmeDescription = "Append licensing information to README."

var README = Subcommand{
	Tag:         "seller",
	Description: readmeDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("readme", flag.ExitOnError)
		silent := Silent(flagSet)
		flagSet.Usage = readmeUsage
		flagSet.Parse(args)
		var existing string
		data, err := ioutil.ReadFile("README.md")
		if err != nil {
			if os.IsNotExist(err) {
				existing = ""
			} else {
				os.Stderr.WriteString("Error reading README.md.\n")
				os.Exit(1)
			}
		} else {
			existing = string(data)
		}
		projectIDs, termsIDs, err := readEntries(paths.CWD)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
		if len(existing) > 0 {
			existing = existing + "\n\n"
		}
		existing = existing + "# Licensing"
		if len(projectIDs) == 0 {
			os.Stderr.WriteString("No License Zero project metadata in package.json.\n")
			os.Exit(1)
		}
		haveReciprocal := false
		haveNoncommercial := false
		for _, terms := range termsIDs {
			if terms == "noncommercial" {
				haveNoncommercial = true
			} else if terms == "reciprocal" {
				haveReciprocal = true
			}
		}
		var summary string
		var availability string
		if haveReciprocal && haveNoncommercial {
			summary = "" +
				"Some contributions to this package " +
				"are free to use in open source under the terms of " +
				"the [License Zero Reciprocal Public License](./LICENSE), " +
				"and some contributions to this package " +
				"are free to use for nomcommercial purposes under the terms of " +
				"the [License Zero Noncommercial Public License](./LICENSE), "
			availability = "" +
				"Licenses for use in closed software, and for long-term commercial use " +
				"are available via [licensezero.com](https://licensezero.com)."
		} else if haveReciprocal {
			summary = "" +
				"This package " +
				"is free to use in open source under the terms of " +
				"the [License Zero Reciprocal Public License](./LICENSE)."
			availability = "" +
				"Licenses for use in closed software " +
				"are available via [licensezero.com](https://licensezero.com)."
		} else if haveNoncommercial {
			summary = "" +
				"This package " +
				"is free to use for noncommercial purposes under the terms of " +
				"the [License Zero Noncommercial Public License](./LICENSE)."
			availability = "" +
				"Licenses for long-term commercial use " +
				"are available via [licensezero.com](https://licensezero.com)."
		} else {
			os.Stderr.WriteString("Unrecognized License Zero project terms.\n")
			os.Exit(1)
		}
		existing = existing + "\n\n" + summary
		existing = existing + "\n\n" + availability
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
			os.Stderr.WriteString("Error writing README.md.\n")
			os.Exit(1)
		}
		if !*silent {
			os.Stdout.WriteString("Wrote to README.md\n")
		}
		os.Exit(0)
	},
}

func readEntries(directory string) ([]string, []string, error) {
	data, err := ioutil.ReadFile("package.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.New("Could not read package.json.")
		} else {
			return nil, nil, err
		}
	}
	var existingMetadata struct {
		LicenseZero []struct {
			License struct {
				ProjectID string `json:"projectID"`
				Terms     string `json:"terms"`
			} `json:"license"`
			AgentSignature    string `json:"agentSignature"`
			LicensorSignature string `json:"licensorSignature"`
		} `json:"licensezero"`
	}
	err = json.Unmarshal(data, &existingMetadata)
	if err != nil {
		return nil, nil, errors.New("Could not parse package.json metadata.")
	}
	// TODO: Validate package.json metadata entries.
	var projectIDs []string
	var terms []string
	for _, entry := range existingMetadata.LicenseZero {
		projectIDs = append(projectIDs, entry.License.ProjectID)
		terms = append(terms, entry.License.Terms)
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
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
