package subcommands

import "encoding/json"
import "flag"
import "fmt"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/inventory"
import "io/ioutil"
import "os"
import "strconv"

const quoteDescription = "Quote missing private licenses."

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
			Fail("Could not read dependeny tree.")
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
				Fail("Error serializing output.")
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
		if len(unlicensed) == 0 {
			os.Exit(0)
		}
		var projectIDs []string
		for _, project := range unlicensed {
			projectIDs = append(projectIDs, project.Envelope.Manifest.ProjectID)
		}
		response, err := api.Quote(projectIDs)
		if err != nil {
			Fail("Error requesting quote.")
		}
		var total uint
		for _, project := range response.Projects {
			var prior *inventory.Project
			for _, candidate := range unlicensed {
				if candidate.Envelope.Manifest.ProjectID == project.ProjectID {
					prior = &candidate
					break
				}
			}
			total += project.Pricing.Private
			fmt.Println("\n- Project: " + project.ProjectID)
			fmt.Println("  Description: " + project.Description)
			fmt.Println("  Repository: " + project.Repository)
			if prior != nil {
				if prior.Envelope.Manifest.Terms == "noncommercial" {
					fmt.Println("  Terms: Noncommercial")
				} else if prior.Envelope.Manifest.Terms == "reciprocal" {
					fmt.Println("  Terms: Reciprocal")
				} else if prior.Envelope.Manifest.Terms == "parity" {
					fmt.Println("  Terms: Parity")
				} else if prior.Envelope.Manifest.Terms == "prosperity" {
					fmt.Println("  Terms: Prosperity")
				}
			}
			fmt.Println("  Licensor: " + project.Licensor.Name + " [" + project.Licensor.Jurisdiction + "]")
			if project.Retracted {
				fmt.Println("  Retracted!")
			}
			if prior != nil {
				if prior.Type != "" {
					fmt.Println("  Type: " + prior.Type)
				}
				if prior.Path != "" {
					fmt.Println("  Path: " + prior.Path)
				}
				if prior.Scope != "" {
					fmt.Println("  Scope: " + prior.Scope)
				}
				if prior.Name != "" {
					fmt.Println("  Name: " + prior.Name)
				}
				if prior.Version != "" {
					fmt.Println("  Version: " + prior.Version)
				}
			}
			fmt.Println("  Price: " + currency(project.Pricing.Private))
			fmt.Printf("\nTotal: %s\n", currency(total))
		}
		if len(unlicensed) == 0 {
			os.Exit(0)
		} else {
			os.Exit(1)
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
	Fail(usage)
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
