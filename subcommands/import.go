package subcommands

import "encoding/json"
import "flag"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "os"

const importDescription = "Import a private license or waiver from a file."

var Import = Subcommand{
	Description: importDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("import", flag.ExitOnError)
		filePath := flagSet.String("file", "", "")
		silent := Silent(flagSet)
		flagSet.Usage = importUsage
		flagSet.Parse(args)
		if *filePath == "" {
			importUsage()
		}
		bytes, err := ioutil.ReadFile(*filePath)
		var initialParse interface{}
		err = json.Unmarshal(bytes, initialParse)
		if err != nil {
			os.Stderr.WriteString("Invalid JSON\n")
			os.Exit(1)
		}
		itemsMap := initialParse.(map[string]interface{})
		if _, ok := itemsMap["license"]; ok {
			license, err := data.ReadLicense(*filePath)
			if err != nil {
				os.Stderr.WriteString("Error reading license\n")
				os.Exit(1)
			}
			// TODO: Validate licenses.
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
			// TODO: Validate waivers.
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
		"  --file FILE  License or waiver file to import.\n" +
		"  --silent     " + silentLine + "\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
