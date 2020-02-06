package api

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
)

// Licensor contains data about a licensor from a vendor API.
type Licensor struct {
	Jurisdiction string
	Name         string
}

func parseLicensor(unstructured interface{}) (l *Licensor, err error) {
	if !validateLicensor(unstructured) {
		return nil, errors.New("invalid licensor")
	}
	err = mapstructure.Decode(unstructured, &l)
	return
}

var licensorValidator *gojsonschema.Schema = nil

func validateLicensor(unstructured interface{}) bool {
	if licensorValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Licensor),
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
