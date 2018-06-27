package subcommands

import "encoding/json"
import "flag"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "net/http"
import "os"
import "strconv"

const purchasedDescription = "Download a bundle of licenses you bought from licensezero.com."

var Purchased = Subcommand{
	Tag:         "buyer",
	Description: purchasedDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("purchased", flag.ExitOnError)
		bundle := flagSet.String("bundle", "", "")
		silent := Silent(flagSet)
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = purchasedUsage
		flagSet.Parse(args)
		if *bundle == "" {
			purchasedUsage()
		}
		response, err := http.Get(*bundle)
		defer response.Body.Close()
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			os.Stderr.WriteString("Error reading " + *bundle + ".\n")
			os.Exit(1)
		}
		var parsed struct {
			Licenses []data.LicenseFile `json:"licenses"`
		}
		err = json.Unmarshal(responseBody, &parsed)
		if err != nil {
			os.Stderr.WriteString("Error parsing license bundle.\n")
			os.Exit(1)
		}
		imported := 0
		for _, license := range parsed.Licenses {
			envelope, err := data.LicenseFileToEnvelope(&license)
			if err != nil {
				os.Stderr.WriteString("Error parsing license for project ID" + license.ProjectID + ".\n")
				continue
			}
			// TODO: Validate licenses.
			err = data.WriteLicense(paths.Home, envelope)
			if err != nil {
				os.Stderr.WriteString("Error writing license for project ID" + license.ProjectID + ".\n")
				continue
			}
			imported++
		}
		if !*silent {
			os.Stdout.WriteString("Imported " + strconv.Itoa(imported) + " licenses.\n")
		}
		os.Exit(0)
	},
}

func purchasedUsage() {
	usage := purchasedDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero purchased --bundle URL\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"bundle URL": "URL of purchase bundle to import.",
			"silent":     silentLine,
		})
	os.Stderr.WriteString(usage)
	os.Exit(1)
}
