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

const licenseDescription = "Write LICENSE and package metadata for your project."

var License = Subcommand{
	Tag:         "seller",
	Description: licenseDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("license", flag.ExitOnError)
		projectID := ProjectID(flagSet)
		noncommercial := flagSet.Bool("noncommercial", false, "Use noncommercial public license.")
		reciprocal := flagSet.Bool("reciprocal", false, "Use reciprocal public license.")
		stack := flagSet.Bool("stack", false, "Stack licensing metadata.")
		silent := Silent(flagSet)
		flagSet.Usage = licenseUsage
		flagSet.Parse(args)
		if *noncommercial && *reciprocal {
			licenseUsage()
		}
		if !*noncommercial && !*reciprocal {
			licenseUsage()
		}
		if *projectID == "" {
			licenseUsage()
		}
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			os.Stderr.WriteString(licensorHint + "\n")
			os.Exit(1)
		}
		var terms string
		if *noncommercial {
			terms = "noncommercial"
		}
		if *reciprocal {
			terms = "reciprocal"
		}
		response, err := api.Public(licensor, *projectID, terms)
		if err != nil {
			os.Stderr.WriteString("Error sending license information request.\n")
			os.Exit(1)
		}
		// Add metadata to package.json.
		newEntry := response.Metadata.LicenseZero
		package_json := path.Join(paths.CWD, "package.json")
		data, err := ioutil.ReadFile(package_json)
		if err != nil {
			os.Stderr.WriteString("Could not read package.json.\n")
			os.Exit(1)
		}
		var existingMetadata interface{}
		err = json.Unmarshal(data, &existingMetadata)
		if err != nil {
			os.Stderr.WriteString("Error parsing package.json.\n")
			os.Exit(1)
		}
		itemsMap := existingMetadata.(map[string]interface{})
		var entries []interface{}
		if _, ok := itemsMap["licensezero"]; ok {
			if entries, ok := itemsMap["licensezero"].([]interface{}); ok {
				if *stack {
					entries = append(entries, newEntry)
				} else {
					os.Stderr.WriteString("package.json already has License Zero metadata.\nUse --stack to stack metadata.\n")
					os.Exit(1)
				}
			} else {
				os.Stderr.WriteString("package.json has an invalid licensezero property.\n")
				os.Exit(1)
			}
		} else {
			if *stack {
				os.Stderr.WriteString("Cannot stack License Zero metadata. There is no preexisting metadata.\n")
				os.Exit(1)
			} else {
				entries = []interface{}{newEntry}
			}
		}
		itemsMap["licensezero"] = entries
		serialized, err := json.Marshal(existingMetadata)
		if err != nil {
			os.Stderr.WriteString("Error serializing new JSON\n")
			os.Exit(1)
		}
		indented := bytes.NewBuffer([]byte{})
		err = json.Indent(indented, serialized, "", "  ")
		if err != nil {
			os.Stderr.WriteString("Error indenting new JSON.\n")
			os.Exit(1)
		}
		err = ioutil.WriteFile(package_json, indented.Bytes(), 0644)
		if !*silent {
			os.Stdout.WriteString("Added metadata to package.json.\n")
		}
		// Append to LICENSE.
		err = writeLICENSE(&response)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
		if !*silent {
			os.Stdout.WriteString("Appended terms to LICENSE.")
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
			return errors.New("Could not open LICENSE.")
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
		signatureLines(response.License.AgentSignature)
	err = ioutil.WriteFile("LICENSE", []byte(toWrite), 0644)
	if err != nil {
		return errors.New("Error writing LICENSE")
	}
	return nil
}

func signatureLines(signature string) string {
	return "" +
		signature[0:32] + "\n" +
		signature[32:64] + "\n" +
		signature[64:96] + "\n" +
		signature[96:]
}

func licenseUsage() {
	usage := licenseDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero license --project ID (--noncommercial | --reciprocal) [--stack]\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"project":       projectIDLine,
			"noncommerical": "Use the noncommercial license.",
			"reciprocal":    "Use the reciprocal license.",
			"silent":        silentLine,
			"stack":         "Stack multiple project metadata entries.",
		})
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
