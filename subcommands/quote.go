package subcommands

import "bytes"
import "encoding/json"
import "flag"
import "fmt"
import "io/ioutil"
import "net/http"
import "os"

type QuoteRequest struct {
	Action   string   `json:"action"`
	Projects []string `json:"projects"`
}

type QuoteResponse struct {
	Projects []QuoteProject `json:"projects"`
}

type QuoteProject struct {
	Licensor    LicensorInformation `json:"licensor"`
	ProjectID   string              `json:"projectID"`
	Description string              `json:"description"`
	Repository  string              `json:"homepage"`
	Pricing     Pricing             `json:"pricing"`
	Retracted   bool                `json:"retracted"`
}

type LicensorInformation struct {
	Name         string
	Jurisdiction string
}

type Pricing struct {
	Private int
}

var Quote = Subcommand{
	Description: "Quote missing private licenses.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("quote", flag.ExitOnError)
		noNoncommercial := flagSet.Bool("no-noncommercial", false, "Ignore L0-NC dependencies.")
		noReciprocal := flagSet.Bool("no-reciprocal", false, "Ignore L0-R dependencies.")
		flagSet.Parse(args)
		projects, err := Inventory(paths, *noNoncommercial, *noReciprocal)
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
				projectIDs = append(projectIDs, project.Manifest.ProjectID)
			}
			bodyData := QuoteRequest{
				Action:   "quote",
				Projects: projectIDs,
			}
			body, err := json.Marshal(bodyData)
			if err != nil {
				os.Stderr.WriteString("Could not construct quote request.")
				os.Exit(1)
			}
			response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
			defer response.Body.Close()
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				os.Stderr.WriteString("Invalid server response.")
				os.Exit(1)
			}
			var parsed QuoteResponse
			json.Unmarshal(responseBody, &parsed)
			total := 0
			for _, project := range parsed.Projects {
				total += project.Pricing.Private
				fmt.Println("\n- Project: " + project.ProjectID)
				fmt.Println("  Description: " + project.Description)
				fmt.Println("  Repository: " + project.Repository)
				// TODO: Terms
				for _, prior := range unlicensed {
					if prior.Manifest.ProjectID == project.ProjectID {
						if prior.Manifest.Terms == "noncommercial" {
							fmt.Println("  Terms: Noncommercial " + prior.Version)
						} else if prior.Manifest.Terms == "reciprocal" {
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
