package api

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
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
