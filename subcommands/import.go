package subcommands

import "encoding/json"
import "flag"
import "github.com/licensezero/cli/data"
import "github.com/licensezero/cli/api"
import "io/ioutil"
import "os"

const importDescription = "Import a private license or waiver from a file."

var Import = Subcommand{
	Tag:         "buyer",
	Description: importDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("import", flag.ExitOnError)
		filePath := flagSet.String("file", "", "")
		silent := Silent(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = importUsage
		flagSet.Parse(args)
		if *filePath == "" {
			importUsage()
		}
		bytes, err := ioutil.ReadFile(*filePath)
		if err != nil {
			os.Stderr.WriteString("Could not read file.")
			os.Exit(1)
		}
		var documentPreview struct {
			Manifest string `json:"manifest"`
		}
		err = json.Unmarshal(bytes, &documentPreview)
		if err != nil {
			os.Stderr.WriteString("Invalid JSON\n")
			os.Exit(1)
		}
		var manifestPreview struct {
			Form string `json:"FORM"`
		}
		err = json.Unmarshal([]byte(documentPreview.Manifest), &manifestPreview)
		if err != nil {
			os.Stderr.WriteString("Invalid manifest JSON\n")
			os.Exit(1)
		}
		os.Stderr.WriteString(manifestPreview.Form + "\n")
		if manifestPreview.Form == "private license" {
			license, err := data.ReadLicense(*filePath)
			if err != nil {
				os.Stderr.WriteString("Error reading license\n")
				os.Exit(1)
			}
			projectResponse, err := api.Project(license.Manifest.Project.ProjectID)
			if err != nil {
				os.Stderr.WriteString("Error fetching project developer information.\n")
				os.Exit(1)
			}
			err = data.CheckLicenseSignature(license, projectResponse.Licensor.PublicKey)
			if err != nil {
				os.Stderr.WriteString("Invalid license signature.")
				os.Exit(1)
			}
			err = data.WriteLicense(paths.Home, license)
			if err != nil {
				os.Stderr.WriteString("Error writing license file.\n")
				os.Exit(1)
			}
		} else {
			waiver, err := data.ReadWaiver(*filePath)
			if err != nil {
				os.Stderr.WriteString("Error reading waiver.\n")
				os.Exit(1)
			}
			projectResponse, err := api.Project(waiver.Manifest.Project.ProjectID)
			if err != nil {
				os.Stderr.WriteString("Error fetching project developer information.")
				os.Exit(1)
			}
			err = data.CheckWaiverSignature(waiver, projectResponse.Licensor.PublicKey)
			if err != nil {
				os.Stderr.WriteString("Invalid waiver signature.")
				os.Exit(1)
			}
			err = data.WriteWaiver(paths.Home, waiver)
			if err != nil {
				os.Stderr.WriteString("Error writing waiver file.\n")
				os.Exit(1)
			}
		}
		if !*silent {
			os.Stdout.WriteString("Imported.\n")
		}
		os.Exit(0)
	},
}

func importUsage() {
	usage := importDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero import --file FILE\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"file FILE": "License or waiver file to import.",
			"silent":    silentLine,
		})
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
