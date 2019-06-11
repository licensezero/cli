package subcommands

import "encoding/json"
import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "net/http"
import "os"
import "strconv"

const importDescription = "Import private licenses and waivers."

// Import saves licenses and waivers to the config directory.
var Import = &Subcommand{
	Tag:         "buyer",
	Description: importDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("import", flag.ExitOnError)
		filePath := flagSet.String("file", "", "")
		bundle := flagSet.String("bundle", "", "")
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = importUsage
		flagSet.Parse(args)
		if *filePath == "" && *bundle == "" {
			importUsage()
		} else if *filePath != "" {
			importFile(paths, filePath, silent)
		} else {
			importBundle(paths, bundle, silent)
		}
	},
}

func importBundle(paths Paths, bundle *string, silent *bool) {
	response, err := http.Get(*bundle)
	if err != nil {
		Fail("Error getting bundle.")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Fail("Error reading " + *bundle + ".")
	}
	var parsed struct {
		Licenses []data.LicenseFile `json:"licenses"`
	}
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		Fail("Error parsing license bundle.")
	}
	imported := 0
	for _, license := range parsed.Licenses {
		envelope, err := data.LicenseFileToEnvelope(&license)
		if err != nil {
			os.Stderr.WriteString("Error parsing license for project ID" + license.ProjectID + ".\n")
			continue
		}
		projectID := envelope.Manifest.Project.ProjectID
		licensor, err := api.Project(projectID)
		if err != nil {
			os.Stderr.WriteString("Error fetching project developer information for " + projectID + ": " + err.Error() + "\n")
			continue
		}
		err = data.CheckLicenseSignature(envelope, licensor.PublicKey)
		if err != nil {
			os.Stderr.WriteString("Invalid license signature for project " + projectID + ".\n")
			continue
		}
		err = data.WriteLicense(paths.Home, envelope)
		if err != nil {
			os.Stderr.WriteString("Error writing license for project ID" + license.ProjectID + ".\n")
			continue
		}
		imported++
	}
	if !*silent {
		os.Stdout.WriteString("Imported " + strconv.Itoa(imported) + " licenses.\n")
	}
	os.Exit(0)
}

func importFile(paths Paths, filePath *string, silent *bool) {
	bytes, err := ioutil.ReadFile(*filePath)
	if err != nil {
		Fail("Could not read file.")
	}
	var documentPreview struct {
		Manifest string `json:"manifest"`
	}
	err = json.Unmarshal(bytes, &documentPreview)
	if err != nil {
		Fail("Invalid JSON")
	}
	var manifestPreview struct {
		Form string `json:"FORM"`
	}
	err = json.Unmarshal([]byte(documentPreview.Manifest), &manifestPreview)
	if err != nil {
		Fail("Invalid manifest JSON")
	}
	if manifestPreview.Form == "private license" {
		license, err := data.ReadLicense(*filePath)
		if err != nil {
			Fail("Error reading license.")
		}
		licensor, err := api.Project(license.Manifest.Project.ProjectID)
		if err != nil {
			Fail("Error fetching project developer information: " + err.Error())
		}
		err = data.CheckLicenseSignature(license, licensor.PublicKey)
		if err != nil {
			Fail("Invalid license signature.")
		}
		err = data.WriteLicense(paths.Home, license)
		if err != nil {
			Fail("Error writing license file.")
		}
	} else {
		waiver, err := data.ReadWaiver(*filePath)
		if err != nil {
			Fail("Error reading waiver.")
		}
		licensor, err := api.Project(waiver.Manifest.Project.ProjectID)
		if err != nil {
			Fail("Error fetching project developer information: " + err.Error())
		}
		err = data.CheckWaiverSignature(waiver, licensor.PublicKey)
		if err != nil {
			Fail("Invalid waiver signature.")
		}
		err = data.WriteWaiver(paths.Home, waiver)
		if err != nil {
			Fail("Error writing waiver file.")
		}
	}
	if !*silent {
		os.Stdout.WriteString("Imported.\n")
	}
	os.Exit(0)
}

func importUsage() {
	usage := importDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero import (--file FILE | --bundle URL)\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"bundle URL": "URL of purchased license bundle.",
			"file FILE":  "License or waiver file to import.",
			"silent":     silentLine,
		})
	Fail(usage)
}
