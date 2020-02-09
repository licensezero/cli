package api

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
)

// Offer represents an offer to sell licenses.
type Offer struct {
	URL      string  `json:"url"`
	SellerID string  `json:"sellerID"`
	Pricing  Pricing `json:"pricing"`
}

// Pricing represents a price list.
type Pricing struct {
	Single    Price `json:"single"`
	Relicense Price `json:"relicense"`
}

// ErrInvalidOffer indicates that an Offer does not conform
// to the JSON schema for offer records.
var ErrInvalidOffer = errors.New("invalid offer")

var offerValidator *gojsonschema.Schema = nil

func (offer *Offer) Validate() error {
	if offerValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Offer),
		)
		if err != nil {
			panic(err)
		}
		offerValidator = schema
	}
	marshaled, err := json.Marshal(offer)
	if err != nil {
		return err
	}
	dataLoader := gojsonschema.NewBytesLoader(marshaled)
	result, err := offerValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return ErrInvalidOffer
	}
	return nil
}
