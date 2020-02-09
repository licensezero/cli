package api

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
)

// Seller represents the party selling licenses.
// Seller is usually the developer of the software being licensed.
type Seller struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
}

// ErrInvalidSeller indicates that a Seller does not conform
// to the JSON schema for seller records.
var ErrInvalidSeller = errors.New("invalid seller")

var sellerValidator *gojsonschema.Schema

// Validate verifies that the Seller conforms
// to the JSON schema for Seller records.
func (seller *Seller) Validate() error {
	if sellerValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Seller),
		)
		if err != nil {
			panic(err)
		}
		sellerValidator = schema
	}
	marshaled, err := json.Marshal(seller)
	if err != nil {
		return err
	}
	dataLoader := gojsonschema.NewBytesLoader(marshaled)
	result, err := sellerValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return ErrInvalidSeller
	}
	return nil
}
