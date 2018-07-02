package subcommands

import "bytes"
import "encoding/json"
import "errors"
import "flag"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"
import "path"

const licenseDescription = "Write license terms and metadata."

// License writes LICENSE and licensezero.json.
var License = Subcommand{
	Tag:         "seller",
	Description: licenseDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("license", flag.ExitOnError)
		projectID := projectIDFlag(flagSet)
		prosperity := flagSet.Bool("prosperity", false, "Use The Prosperity Public License")
		parity := flagSet.Bool("parity", false, "Use The Parity Public License.")
		stack := flagSet.Bool("stack", false, "Stack licensing metadata.")
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = licenseUsage
		flagSet.Parse(args)
		if *prosperity && *parity {
			licenseUsage()
		}
		if !*prosperity && !*parity {
			licenseUsage()
		}
		if *projectID == "" {
			licenseUsage()
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		var terms string
		if *prosperity {
			terms = "prosperity"
		}
		if *parity {
			terms = "parity"
		}
		response, err := api.Public(licensor, *projectID, terms)
		if err != nil {
			Fail("Error sending license information request.")
		}
		checkForLegacyPackageJSON(paths.CWD)
		// Add metadata to licensezero.json.
		type FileStructure struct {
			Version     string        `json:"version"`
			LicenseZero []interface{} `json:"licensezero"`
		}
		var newMetadata FileStructure
		newEntry := response.Metadata.LicenseZero
		licensezeroJSON := path.Join(paths.CWD, "licensezero.json")
		data, err := ioutil.ReadFile(licensezeroJSON)
		if err != nil {
			if os.IsNotExist(err) {
				newMetadata.LicenseZero = []interface{}{newEntry}
			} else {
				Fail("Could not read licensezero.json.")
			}
		} else {
			var existingMetadata FileStructure
			err = json.Unmarshal(data, &existingMetadata)
			if err != nil {
				Fail("Error parsing licensezero.json.")
			}
			entries := existingMetadata.LicenseZero
			if len(existingMetadata.LicenseZero) != 0 {
				if *stack {
					// Check if project already listed.
					for _, entry := range entries {
						if itemsMap, ok := entry.(map[string]interface{}); ok {
							if license, ok := itemsMap["license"].(map[string]interface{}); ok {
								if otherID, ok := license["projectID"].(string); ok {
									if otherID == *projectID {
										Fail("Project ID " + *projectID + " already appears in licensezero.json.")
									}
								}
							}
						}
					}
					entries = append(existingMetadata.LicenseZero, newEntry)
				} else {
					Fail("licensezero.json already has License Zero metadata.\nUse --stack to stack metadata.")
				}
			} else {
				if *stack {
					Fail("Cannot stack License Zero metadata. There is no preexisting metadata.")
				} else {
					entries = []interface{}{newEntry}
				}
			}
			newMetadata.Version = existingMetadata.Version
			newMetadata.LicenseZero = entries
		}
		serialized := new(bytes.Buffer)
		encoder := json.NewEncoder(serialized)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(newMetadata)
		if err != nil {
			Fail("Error serializing new JSON.")
		}
		err = ioutil.WriteFile(licensezeroJSON, serialized.Bytes(), 0644)
		if !*silent {
			os.Stdout.WriteString("Added metadata to licensezero.json.\n")
		}
		// Append to LICENSE.
		err = writeLICENSE(response)
		if err != nil {
			Fail(err.Error())
		}
		if !*silent {
			os.Stdout.WriteString("Appended terms to LICENSE.\n")
		}
		if !*silent {
			os.Stdout.WriteString(
				"" +
					"Make sure to configure your package manager to include licensezero.json\n" +
					"in your package distribution.\n\n" +
					"npm:    Add licensezero.json to the files array of your npm package's\n" +
					"        package.json file, if you have one.\n\n" +
					"Python: Add licensezero.json to MANIFEST.in.\n",
			)
		}
		os.Exit(0)
	},
}

func writeLICENSE(response *api.PublicResponse) error {
	var toWrite string
	existing, err := ioutil.ReadFile("LICENSE")
	if err != nil {
		if os.IsNotExist(err) {
			toWrite = ""
		} else {
			return errors.New("could not open LICENSE")
		}
	} else {
		toWrite = string(existing)
	}
	if len(toWrite) != 0 {
		toWrite = toWrite + "\n\n"
	}
	toWrite = toWrite +
		response.License.Document + "\n\n" +
		"---\n\n" +
		"Licensor Signature (Ed25519):\n\n" +
		signatureLines(response.License.LicensorSignature) + "\n\n" +
		"---\n\n" +
		"Agent Signature (Ed25519):\n\n" +
		signatureLines(response.License.AgentSignature) + "\n"
	err = ioutil.WriteFile("LICENSE", []byte(toWrite), 0644)
	if err != nil {
		return errors.New("Error writing LICENSE")
	}
	return nil
}

func signatureLines(signature string) string {
	return signature[0:64] + "\n" + signature[64:]
}

// Earlier versions of `licensezero` wrote License Zero licensing
// metadata to `licensezero` array properties of `package.json` files
// for npm projects, rather than to separate `licenserzero.json` files.
// If we see a `package.json` file with a `licensezero` property, warn
// the user and instruct them to upgrade.
func checkForLegacyPackageJSON(directoryPath string) {
	packageJSON := path.Join(directoryPath, "package.json")
	data, err := ioutil.ReadFile(packageJSON)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		Fail("Error reading package.json.")
	}
	var packageData struct {
		LegacyMetadata []interface{} `json:"licensezero"`
	}
	err = json.Unmarshal(data, &packageData)
	if err != nil {
		Fail("Error parsing package.json.")
	}
	if len(packageData.LegacyMetadata) != 0 {
		Fail(
			"" +
				"The `licensezero` property in `package.json` is deprecated\n" +
				"in favor of `licensezero.json`.\n" +
				"Remove the `licensezero` property from `package.json`\n" +
				"and run `licensezero license` again.\n",
		)
	}
}

func licenseUsage() {
	usage := licenseDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero license --project ID (--parity | --prosperity) [--stack]\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"project":    projectIDLine,
			"prosperity": "Use the Prosperity Public License.",
			"parity":     "Use The Parity Publice License.",
			"silent":     silentLine,
			"stack":      "Stack multiple project metadata entries.",
		})
	Fail(usage)
}
