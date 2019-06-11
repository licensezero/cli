package api

import "bytes"
import "encoding/json"
import "errors"
import "licensezero.com/cli/data"
import "io/ioutil"
import "net/http"
import "strings"

// AgencyReference includes the text required in agency terms agreement statements to the API.
const AgencyReference = "the agency terms at https://licensezero.com/terms/agency"
const agencyStatement = "I agree to " + AgencyReference + "."

type offerRequest struct {
	Action      string  `json:"action"`
	LicensorID  string  `json:"licensorID"`
	Token       string  `json:"token"`
	Homepage    string  `json:"homepage"`
	Pricing     Pricing `json:"pricing"`
	Description string  `json:"description"`
	Terms       string  `json:"terms"`
}

type offerResponse struct {
	Error     interface{} `json:"error"`
	ProjectID string      `json:"projectID"`
}

// Offer sends an offer API request.
func Offer(licensor *data.Licensor, homepage, description string, private, relicense uint) (string, error) {
	if !strings.HasPrefix(homepage, "https://") && !strings.HasPrefix(homepage, "http://") {
		homepage = "http://" + homepage
	}
	bodyData := offerRequest{
		Action:      "offer",
		LicensorID:  licensor.LicensorID,
		Token:       licensor.Token,
		Description: description,
		Homepage:    homepage,
		Pricing: Pricing{
			Private:   private,
			Relicense: relicense,
		},
		Terms: agencyStatement,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return "", err
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("error sending request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var parsed offerResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return "", err
	}
	if message, ok := parsed.Error.(string); ok {
		return "", errors.New(message)
	}
	return parsed.ProjectID, nil
}
