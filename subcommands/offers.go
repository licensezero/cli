package subcommands

import "encoding/json"
import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const projectsDescription = "List your projects."

// Offers prints the licensor's projects.
var Offers = &Subcommand{
	Tag:         "seller",
	Description: projectsDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("projects", flag.ExitOnError)
		retracted := flagSet.Bool("include-retracted", false, "")
		outputJSON := flagSet.Bool("json", false, "")
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = projectsUsage
		flagSet.Parse(args)
		licensor, err := data.ReadLicensor(paths.Home)
		if err != nil {
			Fail(licensorHint)
		}
		_, projects, err := api.Licensor(licensor.LicensorID)
		if err != nil {
			Fail("Could not fetch licensor information: " + err.Error())
		}
		var filtered []api.OfferInformation
		if *retracted {
			filtered = projects
		} else {
			for _, project := range projects {
				if project.Retracted == "" {
					filtered = append(filtered, project)
				}
			}
		}
		if *outputJSON {
			marshalled, err := json.Marshal(filtered)
			if err != nil {
				Fail("Error serializing output.")
			}
			os.Stdout.WriteString(string(marshalled) + "\n")
			os.Exit(0)
		}
		for i, project := range filtered {
			if i != 0 {
				os.Stdout.WriteString("\n")
			}
			os.Stdout.WriteString("- Offer ID: " + project.OfferID + "\n")
			os.Stdout.WriteString("  Offered:    " + project.Offered + "\n")
			if project.Retracted != "" {
				os.Stdout.WriteString("  Retracted:  " + project.Offered + "\n")
			}
		}
		os.Exit(0)
	},
}

func projectsUsage() {
	usage := projectsDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero projects\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"json":              "Output JSON.",
			"include-retracted": "List retracted projects.",
		})
	Fail(usage)
}
