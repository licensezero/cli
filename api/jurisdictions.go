package api

import (
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
)

var jurisdictionValidator *gojsonschema.Schema = nil

// ValidateJurisdiction determiens whether a string is valid jurisdiction code.
func ValidateJurisdiction(j string) bool {
	if jurisdictionValidator == nil {
		schema, err := gojsonschema.NewSchemaLoader().Compile(
			gojsonschema.NewStringLoader(schemas.Jurisdiction),
		)
		if err != nil {
			panic(err)
		}
		jurisdictionValidator = schema
	}
	dataLoader := gojsonschema.NewGoLoader(j)
	result, err := jurisdictionValidator.Validate(dataLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}
