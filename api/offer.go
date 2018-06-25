package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "net/http"

const AgencyReference = "the agency terms at https://licensezero.com/terms/agency"
const agencyStatement = "I agree to " + AgencyReference + "."

type OfferRequest struct {
	Action      string  `json:"action"`
	LicensorID  string  `json:"licensorID"`
	Token       string  `json:"token"`
	Homepage    string  `json:"homepage"`
	Pricing     Pricing `json:"pricing"`
	Description string  `json:"description"`
	Terms       string  `json:"terms"`
}

type OfferResponse struct {
	Error     interface{} `json:"error"`
	ProjectID string      `json:"projectID"`
}

func Offer(licensor *data.Licensor, homepage, description string, private, relicense uint) (string, error) {
	bodyData := OfferRequest{
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
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var parsed OfferResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return "", err
	}
	if message, ok := parsed.Error.(string); ok {
		return "", errors.New(message)
	}
	return parsed.ProjectID, nil
}
