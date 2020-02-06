package api

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"licensezero.com/licensezero/schemas"
	"net/http"
)

// Offer represents an offer to sell licenses.
type Offer struct {
	URL        string  `json:"url"`
	LicensorID string  `json:"licensorID"`
	Pricing    Pricing `json:"pricing"`
}

// Pricing represents a price list.
type Pricing struct {
	Single    Price `json:"single"`
	Relicense Price `json:"relicense"`
}

// Price represents a price asked or paid.
type Price struct {
	Amount   uint   `json:"amount"`
	Currency string `json:"currency"`
}

func parseOffer(unstructured interface{}) (o *Offer, err error) {
	if !validateOffer(unstructured) {
		return nil, errors.New("invalid offer")
	}
	err = mapstructure.Decode(unstructured, &o)
	return
}

var offerValidator *gojsonschema.Schema = nil

func validateOffer(unstructured interface{}) bool {
	if offerValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Offer),
		)
		if err != nil {
			panic(err)
		}
		offerValidator = schema
	}
	dataLoader := gojsonschema.NewGoLoader(unstructured)
	result, err := offerValidator.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}

// TODO: Figure out how to mock API responses.

// GetOffer requests offer data from a vendor server.
func GetOffer(api string, offerID string) (offer *Offer, err error) {
	response, err := http.Get(api + "/offers/" + offerID)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var unstructured interface{}
	err = json.Unmarshal(body, &unstructured)
	if err != nil {
		return
	}
	offer, err = parseOffer(unstructured)
	if err != nil {
		return
	}
	return
}
