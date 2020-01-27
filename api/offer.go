package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Offer represents data about an offer from a vendor API.
type Offer struct {
	API        string
	URL        string
	LicensorID string
	Pricing    Pricing
}

// Pricing represents pricing for various kinds of licenses.
type Pricing struct {
	Single    Price
	Relicense Price
	Site      Price
}

// Price presents a specific price in a specific currency.
type Price struct {
	Amount   uint
	Currency string
}

// GetOffer fetches information abourt an offer from a vendor API.
func GetOffer(api string, offerID string) (*Offer, error) {
	response, err := http.Get(api + "/offers/" + offerID)
	if err != nil {
		return nil, errors.New("error sending request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}
	var parsed Offer
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, errors.New("error parsing response body")
	}
	return &parsed, nil
}
