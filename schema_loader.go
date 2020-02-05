package cli

import (
	"github.com/xeipuuv/gojsonschema"
)

func schemaLoader() *gojsonschema.SchemaLoader {
	subschemas := []string{
		jurisdictionSchema,
		keySchema,
		priceSchema,
		signatureSchema,
		timeSchema,
		urlSchema,
	}
	loader := gojsonschema.NewSchemaLoader()
	for _, schema := range subschemas {
		loader.AddSchemas(gojsonschema.NewStringLoader(schema))
	}
	return loader
}
