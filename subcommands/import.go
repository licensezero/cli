package subcommands

import "encoding/json"
import "fmt"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

var Import = Subcommand{
	Description: "Import a private license or waiver from file.",
	Handler: func(args []string, paths Paths) {
		if len(args) != 1 {
			os.Stderr.WriteString(`Import a private license or waiver from file.

Usage:
  <file>
`)
			os.Exit(1)
		} else {
			filePath := args[0]
			fmt.Println(filePath)
			bytes, err := ioutil.ReadFile(filePath)
			var initialParse interface{}
			err = json.Unmarshal(bytes, initialParse)
			if err != nil {
				os.Stderr.WriteString("invalid JSON")
				os.Exit(1)
			}
			itemsMap := initialParse.(map[string]interface{})
			if _, ok := itemsMap["license"]; ok {
				license, err := data.ReadLicense(filePath)
				if err != nil {
					os.Stderr.WriteString("error reading license")
					os.Exit(1)
				}
				// TODO: Validate licenses.
				err = data.WriteLicense(paths.Home, license)
				if err != nil {
					os.Stderr.WriteString("error writing license file")
					os.Exit(1)
				}
			} else {
				waiver, err := data.ReadWaiver(filePath)
				if err != nil {
					os.Stderr.WriteString("error reading waiver")
					os.Exit(1)
				}
				// TODO: Validate waivers.
				err = data.WriteWaiver(paths.Home, waiver)
				if err != nil {
					os.Stderr.WriteString("error writing waiver file")
					os.Exit(1)
				}
			}
			os.Exit(0)
		}
	},
}
