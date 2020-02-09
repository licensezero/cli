package api

import (
	"errors"
	"github.com/mitchellh/mapstructure"
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

func parseSeller(unstructured interface{}) (l *Seller, err error) {
	if !validateSeller(unstructured) {
		return nil, errors.New("invalid seller")
	}
	err = mapstructure.Decode(unstructured, &l)
	return
}

var licensorValidator *gojsonschema.Schema = nil

func validateSeller(unstructured interface{}) bool {
	if licensorValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Seller),
		)
		if err != nil {
			panic(err)
		}
		licensorValidator = schema
	}
	dataLoader := gojsonschema.NewGoLoader(unstructured)
	result, err := licensorValidator.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}
