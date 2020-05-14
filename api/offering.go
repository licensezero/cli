package api

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "net/http"

type offeringRequest struct {
	Action  string `json:"action"`
	OfferID string `json:"offerID"`
}

type offeringResponse struct {
	Error    interface{}         `json:"error"`
	Licensor LicensorInformation `json:"licensor"`
}

// Offering sends an offering API request.
func Offering(offerID string) (*LicensorInformation, error) {
	bodyData := offeringRequest{
		Action:  "offering",
		OfferID: offerID,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, errors.New("error encoding agent key request body")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("error sending request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading agent key response body")
	}
	var parsed offeringResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, errors.New("error parsing agent key response body")
	}
	if message, ok := parsed.Error.(string); ok {
		return nil, errors.New(message)
	}
	return &parsed.Licensor, nil
}
