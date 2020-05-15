package api

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "net/http"

type developerRequest struct {
	Action      string `json:"action"`
	DeveloperID string `json:"developerID"`
}

// OfferInformation describes information on an offer from an API Developer request.
type OfferInformation struct {
	OfferID   string `json:"offerID"`
	Offered   string `json:"offered"`
	Retracted string `json:"retracted,omitempty"`
}

type developerResponse struct {
	Error        interface{}        `json:"error"`
	Name         string             `json:"name"`
	Jurisdiction string             `json:"jurisdiction"`
	PublicKey    string             `json:"publicKey"`
	Offers       []OfferInformation `json:"offers"`
}

// DeveloperInformation describes API data about a developer.
type DeveloperInformation struct {
	Name         string
	Jurisdiction string
	PublicKey    string
}

// Developer sends a developer API request.
func Developer(developerID string) (*DeveloperInformation, []OfferInformation, error) {
	bodyData := developerRequest{
		Action:      "developer",
		DeveloperID: developerID,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, nil, errors.New("error encoding developer request body")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, errors.New("error sending developer request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, errors.New("error reading developer response body")
	}
	var parsed developerResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, nil, errors.New("error parsing developer response body")
	}
	if message, ok := parsed.Error.(string); ok {
		return nil, nil, errors.New(message)
	}
	developer := DeveloperInformation{
		Name:         parsed.Name,
		Jurisdiction: parsed.Jurisdiction,
		PublicKey:    parsed.PublicKey,
	}
	return &developer, parsed.Offers, nil
}
