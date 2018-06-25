package subcommands

import "encoding/json"
import "flag"
import "fmt"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/inventory"
import "io/ioutil"
import "os"
import "strconv"

const quoteDescription = "Quote the cost of private licenses you need."

var Quote = Subcommand{
	Tag:         "buyer",
	Description: quoteDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("quote", flag.ExitOnError)
		noNoncommercial := NoNoncommercial(flagSet)
		noReciprocal := NoReciprocal(flagSet)
		outputJSON := flagSet.Bool("json", false, "")
		flagSet.SetOutput(ioutil.Discard)
		flagSet.Usage = quoteUsage
		flagSet.Parse(args)
		projects, err := inventory.Inventory(paths.Home, paths.CWD, *noNoncommercial, *noReciprocal)
		if err != nil {
			os.Stderr.WriteString("Could not read dependeny tree.\n")
			os.Exit(1)
		}
		licensable := projects.Licensable
		licensed := projects.Licensed
		waived := projects.Waived
		unlicensed := projects.Unlicensed
		ignored := projects.Ignored
		invalid := projects.Invalid
		if *outputJSON {
			marshalled, err := json.Marshal(projects)
			if err != nil {
				os.Stderr.WriteString("Error serializing output.\n")
				os.Exit(1)
			}
			os.Stdout.WriteString(string(marshalled) + "\n")
			os.Exit(0)
		}
		if len(licensable) == 0 {
			fmt.Println("No License Zero dependencies found.")
			os.Exit(0)
		}
		fmt.Printf("License Zero Projects: %d\n", len(licensable))
		fmt.Printf("Licensed: %d\n", len(licensed))
		fmt.Printf("Waived: %d\n", len(waived))
		fmt.Printf("Ignored: %d\n", len(ignored))
		fmt.Printf("Unlicensed: %d\n", len(unlicensed))
		fmt.Printf("Invalid: %d\n", len(invalid))
		var projectIDs []string
		for _, project := range unlicensed {
			projectIDs = append(projectIDs, project.Envelope.Manifest.ProjectID)
		}
		response, err := api.Quote(projectIDs)
		if err != nil {
			os.Stderr.WriteString("Error requesting quote.\n")
			os.Exit(1)
		}
		var total uint
		for _, project := range response.Projects {
			total += project.Pricing.Private
			fmt.Println("\n- Project: " + project.ProjectID)
			fmt.Println("  Description: " + project.Description)
			fmt.Println("  Repository: " + project.Repository)
			for _, prior := range unlicensed {
				if prior.Envelope.Manifest.ProjectID == project.ProjectID {
					if prior.Envelope.Manifest.Terms == "noncommercial" {
						fmt.Println("  Terms: Noncommercial " + prior.Version)
					} else if prior.Envelope.Manifest.Terms == "reciprocal" {
						fmt.Println("  Terms: Reciprocal " + prior.Version)
					}
					break
				}
			}
			fmt.Println("  Licensor: " + project.Licensor.Name + " [" + project.Licensor.Jurisdiction + "]")
			if project.Retracted {
				fmt.Println("  Retracted!")
			}
			fmt.Println("  Price: " + currency(project.Pricing.Private))
			fmt.Printf("\nTotal: %s\n", currency(total))
			os.Exit(0)
		}
	},
}

func quoteUsage() {
	usage := quoteDescription + "\n\n" +
		"Usage:\n" +
		"  licensezero quote\n\n" +
		"Options:\n" +
		flagsList(map[string]string{
			"json":             "Output JSON.",
			"no-noncommercial": noNoncommercialLine,
			"no-reciprocal":    noReciprocalLine,
		})
	os.Stderr.WriteString(usage)
	os.Exit(1)
}

func currency(cents uint) string {
	if cents < 100 {
		if cents < 10 {
			return "$0.0" + strconv.Itoa(int(cents))
		} else {
			return "$0." + strconv.Itoa(int(cents))
		}
	} else {
		asString := fmt.Sprintf("%d", cents)
		return "$" + asString[:len(asString)-2] + "." + asString[len(asString)-2:]
	}
}
