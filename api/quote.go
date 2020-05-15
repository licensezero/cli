package api

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "net/http"

type quoteRequest struct {
	Action string   `json:"action"`
	Offers []string `json:"offers"`
}

// QuoteOffer describes the data the API provides on quoted contribution sets.
type QuoteOffer struct {
	Developer   DeveloperInformation `json:"developer"`
	OfferID     string               `json:"offerID"`
	Description string               `json:"description"`
	Repository  string               `json:"homepage"`
	Pricing     Pricing              `json:"pricing"`
	Retracted   bool                 `json:"retracted"`
}

// DeveloperInformation describes API data about a developer.
type DeveloperInformation struct {
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
func Quote(offerIDs []string) ([]QuoteOffer, error) {
	bodyData := quoteRequest{
		Action: "quote",
		Offers: offerIDs,
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
		Error  interface{}  `json:"error"`
		Offers []QuoteOffer `json:"offers"`
	}
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, err
	}
	if message, ok := parsed.Error.(string); ok {
		return nil, errors.New(message)
	}
	return parsed.Offers, nil
}
