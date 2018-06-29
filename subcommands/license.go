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

var License = Subcommand{
	Tag:         "seller",
	Description: licenseDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("license", flag.ExitOnError)
		projectID := ProjectID(flagSet)
		prosperity := flagSet.Bool("prosperity", false, "Use The Prosperity Public License")
		parity := flagSet.Bool("parity", false, "Use The Parity Public License.")
		stack := flagSet.Bool("stack", false, "Stack licensing metadata.")
		silent := Silent(flagSet)
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
		// Add metadata to package.json.
		newEntry := response.Metadata.LicenseZero
		package_json := path.Join(paths.CWD, "package.json")
		data, err := ioutil.ReadFile(package_json)
		if err != nil {
			Fail("Could not read package.json.")
		}
		var existingMetadata interface{}
		err = json.Unmarshal(data, &existingMetadata)
		if err != nil {
			Fail("Error parsing package.json.")
		}
		itemsMap := existingMetadata.(map[string]interface{})
		var entries []interface{}
		if _, ok := itemsMap["licensezero"]; ok {
			if entries, ok := itemsMap["licensezero"].([]interface{}); ok {
				if *stack {
					entries = append(entries, newEntry)
				} else {
					Fail("package.json already has License Zero metadata.\nUse --stack to stack metadata.")
				}
			} else {
				Fail("package.json has an invalid licensezero property.")
			}
		} else {
			if *stack {
				Fail("Cannot stack License Zero metadata. There is no preexisting metadata.")
			} else {
				entries = []interface{}{newEntry}
			}
		}
		itemsMap["licensezero"] = entries
		serialized, err := json.Marshal(existingMetadata)
		if err != nil {
			Fail("Error serializing new JSON.")
		}
		indented := bytes.NewBuffer([]byte{})
		err = json.Indent(indented, serialized, "", "  ")
		if err != nil {
			Fail("Error indenting new JSON.")
		}
		err = ioutil.WriteFile(package_json, indented.Bytes(), 0644)
		if !*silent {
			os.Stdout.WriteString("Added metadata to package.json.\n")
		}
		// Append to LICENSE.
		err = writeLICENSE(&response)
		if err != nil {
			Fail(err.Error())
		}
		if !*silent {
			os.Stdout.WriteString("Appended terms to LICENSE.\n")
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
	return signature[0:64] + "\n" + signature[64:]
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
