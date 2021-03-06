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

// OfferingResponse represents information about an offer.
type OfferingResponse struct {
	Error       interface{}          `json:"error"`
	Developer   DeveloperInformation `json:"developer"`
	Pricing     Pricing              `json:"pricing"`
	Homepage    string               `json:"homepage"`
	Description string               `json:"description"`
	Lock        LockInformation      `json:"lock"`
	Commission  uint                 `json:"commission"`
}

// LockInformation represents information about pricing locks on offers.
type LockInformation struct {
	Locked string `json:"locked"`
	Unlock string `json:"unlock"`
	Price  uint   `json:"price"`
}

// Offering sends an offering API request.
func Offering(offerID string) (*OfferingResponse, error) {
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
	var parsed OfferingResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, errors.New("error parsing agent key response body")
	}
	if message, ok := parsed.Error.(string); ok {
		return nil, errors.New(message)
	}
	return &parsed, nil
}
