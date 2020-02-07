package schemas

import (
	"github.com/xeipuuv/gojsonschema"
)

//go:generate ./goifyschema Artifact https://schemas.licensezero.com/1.0.0-pre/artifact.json
//go:generate ./goifyschema Bundle https://schemas.licensezero.com/1.0.0-pre/bundle.json
//go:generate ./goifyschema Currency https://schemas.licensezero.com/1.0.0-pre/currency.json
//go:generate ./goifyschema Digest https://schemas.licensezero.com/1.0.0-pre/digest.json
//go:generate ./goifyschema Jurisdiction https://schemas.licensezero.com/1.0.0-pre/jurisdiction.json
//go:generate ./goifyschema Key https://schemas.licensezero.com/1.0.0-pre/key.json
//go:generate ./goifyschema Ledger https://schemas.licensezero.com/1.0.0-pre/ledger.json
//go:generate ./goifyschema Offer https://schemas.licensezero.com/1.0.0-pre/offer.json
//go:generate ./goifyschema Order https://schemas.licensezero.com/1.0.0-pre/order.json
//go:generate ./goifyschema Price https://schemas.licensezero.com/1.0.0-pre/price.json
//go:generate ./goifyschema Receipt https://schemas.licensezero.com/1.0.0-pre/receipt.json
//go:generate ./goifyschema Signature https://schemas.licensezero.com/1.0.0-pre/signature.json
//go:generate ./goifyschema Time https://schemas.licensezero.com/1.0.0-pre/time.json
//go:generate ./goifyschema URL https://schemas.licensezero.com/1.0.0-pre/url.json

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
