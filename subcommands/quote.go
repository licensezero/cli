package subcommands

import "flag"
import "fmt"
import "github.com/licensezero/cli/api"
import "github.com/licensezero/cli/inventory"
import "os"

const quoteDescription = "Quote missing private licenses."

var Quote = Subcommand{
	Description: quoteDescription,
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("quote", flag.ExitOnError)
		noNoncommercial := NoNoncommercial(flagSet)
		noReciprocal := NoReciprocal(flagSet)
		flagSet.Usage = quoteUsage
		flagSet.Parse(args)
		projects, err := inventory.Inventory(paths.Home, paths.CWD, *noNoncommercial, *noReciprocal)
		if err != nil {
			os.Stderr.WriteString("Could not read dependeny tree.")
			os.Exit(1)
		} else {
			licensable := projects.Licensable
			licensed := projects.Licensed
			waived := projects.Waived
			unlicensed := projects.Unlicensed
			ignored := projects.Ignored
			invalid := projects.Invalid
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
				os.Stderr.WriteString("Error requesting quote.")
				os.Exit(1)
			}
			total := 0
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
			}

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
		"  --no-noncommercial  " + noNoncommercialLine + "\n" +
		"  --no-reciprocal     " + noReciprocalLine + "\n" +
		"  --do-not-open       " + doNotOpenLine + "\n"
	os.Stderr.WriteString(usage)
	os.Exit(1)
}

func currency(cents int) string {
	if cents < 100 {
		if cents < 10 {
			return "$0.0" + string(cents)
		} else {
			return "$0." + string(cents)
		}
	} else {
		asString := fmt.Sprintf("%d", cents)
		return "$" + asString[:len(asString)-2] + "." + asString[len(asString)-2:]
	}
}
