package schemas

import (
	"github.com/xeipuuv/gojsonschema"
)

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
