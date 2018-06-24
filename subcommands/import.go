package subcommands

import "encoding/json"
import "flag"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

const importDescription = "Import a private license or waiver from file"

var Import = Subcommand{
	Description: importDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("import", flag.ExitOnError)
		filePath := flagSet.String("file", "", "")
		flagSet.Usage = importUsage
		flagSet.Parse(args)
		if *filePath == "" {
			importUsage()
		}
		bytes, err := ioutil.ReadFile(*filePath)
		var initialParse interface{}
		err = json.Unmarshal(bytes, initialParse)
		if err != nil {
			os.Stderr.WriteString("invalid JSON")
			os.Exit(1)
		}
		itemsMap := initialParse.(map[string]interface{})
		if _, ok := itemsMap["license"]; ok {
			license, err := data.ReadLicense(*filePath)
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
			waiver, err := data.ReadWaiver(*filePath)
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
	},
}

func importUsage() {
	usage := importDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero import --file FILE\n\n" +
		"Options:\n" +
		"  --file FILE  License or waiver file to import.\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
