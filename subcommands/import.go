package subcommands

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
	"net/http"
	"os"
	"strconv"
)

const importDescription = "Import receipts."

// Import saves receipts to disk.
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
		Licenses []user.LicenseFile `json:"licenses"`
	}
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		Fail("Error parsing license bundle.")
	}
	imported := 0
	for _, license := range parsed.Licenses {
		envelope, err := user.LicenseFileToEnvelope(&license)
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
		err = user.CheckLicenseSignature(envelope, licensor.PublicKey)
		if err != nil {
			os.Stderr.WriteString("Invalid license signature for project " + projectID + ".\n")
			continue
		}
		err = user.WriteLicense(paths.Home, envelope)
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
	data, err := ioutil.ReadFile(*filePath)
	if err != nil {
		Fail("Could not read file.")
	}
	var unstructured interface{}
	err = json.Unmarshal(data, &unstructured)
	if err != nil {
		Fail("Invalid JSON.")
	}
	err = user.ValidateReceipt(unstructured)
	if err != nil {
		Fail("Invalid receipt.")
	}
	err = user.ValidateSignature(unstructured)
	if err != nil {
		Fail("Invalid signature.")
	}
	receipt := user.ParseReceipt(unstructured)
	err = ioutil.WriteFile(
		user.ReceiptPath(paths.Home, receipt.API, receipt.OfferID),
		data,
		0700,
	)
	if err != nil {
		Fail("Error writing file.")
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
