package subcommands

import "encoding/json"
import "flag"
import "licensezero.com/cli/api"
import "licensezero.com/cli/data"
import "io/ioutil"
import "os"

const projectsDescription = "List your projects."

// Offers prints the developer's projects.
var Offers = &Subcommand{
	Description: projectsDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("projects", flag.ExitOnError)
		retracted := flagSet.Bool("include-retracted", false, "")
		outputJSON := flagSet.Bool("json", false, "")
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = projectsUsage
		flagSet.Parse(args)
		developer, err := data.ReadDeveloper(paths.Home)
		if err != nil {
			Fail(developerHint)
		}
		_, projects, err := api.Developer(developer.DeveloperID)
		if err != nil {
			Fail("Could not fetch developer information: " + err.Error())
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
		type outputItem struct {
			OfferID     string              `json:"offerID"`
			Offered     string              `json:"offered"`
			Retracted   string              `json:"retracted,omitempty"`
			Homepage    string              `json:"homepage"`
			Description string              `json:"description"`
			Pricing     api.Pricing         `json:"pricing"`
			Lock        api.LockInformation `json:"lock"`
			Commission  uint                `json:"commission"`
		}
		var output []outputItem
		for _, project := range filtered {
			info, err := api.Offering(project.OfferID)
			if err != nil {
				Fail("Error fetching info for offer:" + project.OfferID)
			}
			output = append(output, outputItem{
				OfferID:     project.OfferID,
				Offered:     project.Offered,
				Retracted:   project.Retracted,
				Pricing:     info.Pricing,
				Homepage:    info.Homepage,
				Description: info.Description,
				Lock:        info.Lock,
				Commission:  info.Commission,
			})
		}
		if *outputJSON {
			marshalled, err := json.Marshal(output)
			if err != nil {
				Fail("Error serializing output.")
			}
			os.Stdout.WriteString(string(marshalled) + "\n")
			os.Exit(0)
		}
		for i, item := range output {
			if i != 0 {
				os.Stdout.WriteString("\n")
			}
			os.Stdout.WriteString("- Offer ID: " + item.OfferID + "\n")
			os.Stdout.WriteString("  Offered:  " + item.Offered + "\n")
			if item.Retracted != "" {
				os.Stdout.WriteString("  Retracted:  " + item.Offered + "\n")
			}
			os.Stdout.WriteString("  Homepage: " + item.Homepage + "\n")
			os.Stdout.WriteString("  Description: " + item.Description + "\n")
			os.Stdout.WriteString("  Pricing:\n")
			os.Stdout.WriteString("    Private: " + currency(item.Pricing.Private) + "\n")
			if item.Lock.Locked != "" {
				os.Stdout.WriteString("  Locked:\n")
				os.Stdout.WriteString("    Date:    " + item.Lock.Locked + "\n")
				os.Stdout.WriteString("    Expires: " + item.Lock.Unlock + "\n")
				os.Stdout.WriteString("    Price:   " + currency(item.Lock.Price) + "\n")
			}
			os.Stdout.WriteString("  Commission: " + commission(item.Commission) + "\n")
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
