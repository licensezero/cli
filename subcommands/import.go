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
	var parsed []user.Receipt
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		Fail("Error parsing license bundle.")
	}
	client := api.NewClient(http.DefaultTransport)
	imported := 0
	for _, receipt := range parsed {
		err = receipt.VerifySignature()
		server := receipt.License.Values.API
		offerID := receipt.License.Values.OfferID
		if err != nil {
			os.Stderr.WriteString("Invalid signature for " + server + "/offers/" + offerID + "\n")
			continue
		}
		_, err := client.Offer(server, offerID)
		if err != nil {
			os.Stderr.WriteString("Error fetching offer.\n")
			continue
		}
		err = receipt.Save(paths.Home)
		if err != nil {
			os.Stderr.WriteString("Error saving receipt for " + server + "/offers/" + offerID + "\n")
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
	var receipt user.Receipt
	err = json.Unmarshal(data, &receipt)
	if err != nil {
		Fail("Invalid receipt.")
	}
	err = receipt.VerifySignature()
	if err != nil {
		Fail("Invalid signature.")
	}
	err = receipt.Save(paths.Home)
	if err != nil {
		Fail("Error saving receipt.")
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
