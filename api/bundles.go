package api

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
)

// Bundle represents a bundle of receipts.
type Bundle struct {
	Created  string    `json:"created"`
	Receipts []Receipt `json:"receipts"`
}

// ErrInvalidBundle indicates that an Bundle does not conform
// to the JSON schema for bundles.
var ErrInvalidBundle = errors.New("invalid bundle")

var bundleValidator *gojsonschema.Schema = nil

// Validate checks that a Bundle conforms to the
// JSON schema for bundle data.
func (bundle *Bundle) Validate() error {
	if bundleValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Bundle),
		)
		if err != nil {
			panic(err)
		}
		bundleValidator = schema
	}
	marshaled, err := json.Marshal(bundle)
	if err != nil {
		return err
	}
	dataLoader := gojsonschema.NewBytesLoader(marshaled)
	result, err := bundleValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return ErrInvalidBundle
	}
	return nil
}
