package subcommands

import "bytes"
import "encoding/json"
import "flag"
import "fmt"
import "github.com/skratchdot/open-golang/open"
import "io/ioutil"
import "net/http"
import "os"

type BuyRequest struct {
	Action       string   `json:"action"`
	Projects     []string `json:"projects"`
	Name         string   `json:"licensee"`
	Jurisdiction string   `json:"jurisdiction"`
	EMail        string   `json:"email"`
	Person       string   `json:"person"`
}

type BuyResponse struct {
	Location string `json:"location"`
}

var Buy = Subcommand{
	Description: "Buy missing private licenses.",
	Handler: func(args []string, paths Paths) {
		flagSet := flag.NewFlagSet("buy", flag.ExitOnError)
		doNotOpen := DoNotOpen(flagSet)
		noNoncommercial := flagSet.Bool("no-noncommercial", false, "Ignore L0-NC dependencies.")
		noReciprocal := flagSet.Bool("no-reciprocal", false, "Ignore L0-R dependencies.")
		flagSet.Parse(args)
		identity, err := readIdentity(paths.Home)
		if err != nil {
			os.Stderr.WriteString("Create an identity with `licensezero identify` first.")
			os.Exit(1)
		}
		projects, err := Inventory(paths, *noNoncommercial, *noReciprocal)
		if err != nil {
			os.Stderr.WriteString("Could not read dependeny tree.")
			os.Exit(1)
		} else {
			licensable := projects.Licensable
			unlicensed := projects.Unlicensed
			if len(licensable) == 0 {
				fmt.Println("No License Zero depedencies found.")
				os.Exit(0)
			}
			if len(unlicensed) == 0 {
				fmt.Println("No private licenses to buy.")
				os.Exit(0)
			}
			var projectIDs []string
			for _, project := range unlicensed {
				projectIDs = append(projectIDs, project.Manifest.ProjectID)
			}
			bodyData := BuyRequest{
				Action:       "order",
				Projects:     projectIDs,
				Name:         identity.Name,
				Jurisdiction: identity.Jurisdiction,
				EMail:        identity.EMail,
				Person:       "I am a person, not a legal entity.",
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
				os.Stderr.WriteString("Invalid server response.\n")
				os.Exit(1)
			}
			var parsed BuyResponse
			json.Unmarshal(responseBody, &parsed)
			location := parsed.Location
			url := "https://licensezero.com" + location
			os.Stdout.WriteString(url + "\n")
			if !*doNotOpen {
				open.Run(url)
			}
			os.Exit(0)
		}
	},
}
