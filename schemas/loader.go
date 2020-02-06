package schemas

import (
	"github.com/xeipuuv/gojsonschema"
)

// Loader preloads various subschemas. Preloading the
// subschemas allows us to compile top-level schemas
// without making any network calls for the subschemas they
// reference.
func Loader() *gojsonschema.SchemaLoader {
	subschemas := []string{
		Jurisdiction,
		Key,
		Price,
		Signature,
		Time,
		URL,
	}
	loader := gojsonschema.NewSchemaLoader()
	for _, schema := range subschemas {
		loader.AddSchemas(gojsonschema.NewStringLoader(schema))
	}
	return loader
}
