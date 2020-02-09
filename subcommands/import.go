package subcommands

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
	"net/http"
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
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
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
			return importFile(filePath, silent, stdout, stderr, client)
		} else {
			return importBundle(bundle, silent, stdout, stderr, client)
		}
	},
}

func importBundle(bundle *string, silent *bool, stdout, stderr io.StringWriter, client *http.Client) int {
	response, err := client.Get(*bundle)
	if err != nil {
		stderr.WriteString("Error getting bundle.\n")
		stderr.WriteString(err.Error() + "\n")
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
	imported := 0
	for _, receipt := range parsed.Receipts {
		err = importReceipt(&receipt)
		if err != nil {
			displayError(err, stderr)
			continue
		}
		imported++
	}
	if !*silent {
		stdout.WriteString("Imported " + strconv.Itoa(imported))
		if imported == 1 {
			stdout.WriteString(" license.\n")
		} else {
			stdout.WriteString(" licenses.\n")
		}
	}
	return 0
}

func importFile(filePath *string, silent *bool, stdout, stderr io.StringWriter, client *http.Client) int {
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
	err = importReceipt(&receipt)
	if err != nil {
		displayError(err, stderr)
		return 1
	}
	if !*silent {
		stdout.WriteString("Imported " + *filePath + ".\n")
	}
	return 0
}

func displayError(err error, stderr io.StringWriter) {
	switch {
	case errors.Is(err, api.ErrInvalidReceipt):
		stderr.WriteString("Invalid receipt.")
	case errors.Is(err, api.ErrInvalidSignaure):
		stderr.WriteString("Invalid signature.")
	default:
		stderr.WriteString(err.Error())
	}
	stderr.WriteString("\n")
}

func importReceipt(receipt *api.Receipt) (err error) {
	err = receipt.Validate()
	if err != nil {
		return
	}
	err = receipt.VerifySignature()
	if err != nil {
		return
	}
	err = user.SaveReceipt(receipt)
	return
}
