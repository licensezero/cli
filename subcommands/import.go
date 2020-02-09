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

var importUsage = importDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero import (--file FILE | --bundle URL)\n\n" +
	"Options:\n" +
	flagsList(map[string]string{
		"bundle URL": "URL of purchased license bundle.",
		"file FILE":  "License or waiver file to import.",
		"silent":     silentLine,
	})

// Import saves receipts to disk.
var Import = &Subcommand{
	Tag:         "buyer",
	Description: importDescription,
	Handler: func(args []string, stdin, stdout, stderr *os.File) int {
		flagSet := flag.NewFlagSet("import", flag.ExitOnError)
		filePath := flagSet.String("file", "", "")
		bundle := flagSet.String("bundle", "", "")
		silent := silentFlag(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = func() {
			stderr.WriteString(importUsage)
		}
		err := flagSet.Parse(args)
		if err != nil {
			return 1
		}
		if *filePath == "" && *bundle == "" {
			stderr.WriteString(importUsage)
			return 1
		} else if *filePath != "" {
			return importFile(filePath, silent, stdout, stderr)
		} else {
			return importBundle(bundle, silent, stdout, stderr)
		}
	},
}

func importBundle(bundle *string, silent *bool, stdout, stderr *os.File) int {
	response, err := http.Get(*bundle)
	if err != nil {
		stderr.WriteString("Error getting bundle.\n")
		return 1
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		stderr.WriteString("Error reading " + *bundle + ".\n")
		return 1
	}
	var parsed api.Bundle
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		stderr.WriteString("Error parsing license bundle.\n")
		return 1
	}
	client := api.NewClient(http.DefaultTransport)
	imported := 0
	for _, receipt := range parsed.Receipts {
		err = receipt.VerifySignature()
		server := receipt.License.Values.API
		offerID := receipt.License.Values.OfferID
		if err != nil {
			stderr.WriteString("Invalid signature for " + server + "/offers/" + offerID + "\n")
			continue
		}
		_, err := client.Offer(server, offerID)
		if err != nil {
			stderr.WriteString("Error fetching offer.\n")
			continue
		}
		err = user.SaveReceipt(&receipt)
		if err != nil {
			stderr.WriteString("Error saving receipt for " + server + "/offers/" + offerID + "\n")
			continue
		}
		imported++
	}
	if !*silent {
		stdout.WriteString("Imported " + strconv.Itoa(imported) + " licenses.\n")
	}
	return 0
}

func importFile(filePath *string, silent *bool, stdout, stderr *os.File) int {
	data, err := ioutil.ReadFile(*filePath)
	if err != nil {
		stderr.WriteString("Could not read file.\n")
		return 1
	}
	var receipt api.Receipt
	err = json.Unmarshal(data, &receipt)
	if err != nil {
		stderr.WriteString("Invalid JSON.\n")
		return 1
	}
	err = receipt.Validate()
	if err != nil {
		stderr.WriteString("Invalid receipt.\n")
		return 1
	}
	err = receipt.VerifySignature()
	if err != nil {
		stderr.WriteString("Invalid signature.\n")
		return 1
	}
	err = user.SaveReceipt(&receipt)
	if err != nil {
		stderr.WriteString("Error saving receipt.\n")
		return 1
	}
	if !*silent {
		stdout.WriteString("Imported " + *filePath + ".\n")
	}
	return 0
}
