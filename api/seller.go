package api

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
)

// Seller contains data about a licensor from a vendor API.
type Seller struct {
	Jurisdiction string
	Name         string
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
