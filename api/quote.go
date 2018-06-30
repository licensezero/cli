package api

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "net/http"

type quoteRequest struct {
	Action   string   `json:"action"`
	Projects []string `json:"projects"`
}

// QuoteProject describes the data the API projects on quoted projects.
type QuoteProject struct {
	Licensor    LicensorInformation `json:"licensor"`
	ProjectID   string              `json:"projectID"`
	Description string              `json:"description"`
	Repository  string              `json:"homepage"`
	Pricing     Pricing             `json:"pricing"`
	Retracted   bool                `json:"retracted"`
}

// LicensorInformation describes API data about a licensor.
type LicensorInformation struct {
	Name         string
	Jurisdiction string
	PublicKey    string
}

// Pricing describes private license pricing data.
type Pricing struct {
	Private   uint `json:"private"`
	Relicense uint `json:"relicense,omitempty"`
}

// Quote sends a quote API request.
func Quote(projectIDs []string) ([]QuoteProject, error) {
	bodyData := quoteRequest{
		Action:   "quote",
		Projects: projectIDs,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("error sending request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Error    interface{}    `json:"error"`
		Projects []QuoteProject `json:"projects"`
	}
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, err
	}
	if message, ok := parsed.Error.(string); ok {
		return nil, errors.New(message)
	}
	return parsed.Projects, nil
}
