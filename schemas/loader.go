package schemas

import (
	"github.com/xeipuuv/gojsonschema"
)

//go:generate ./goifyschema Artifact https://protocol.licensezero.com/1.0.0-pre/artifact.json
//go:generate ./goifyschema Broker https://protocol.licensezero.com/1.0.0-pre/broker.json
//go:generate ./goifyschema Bundle https://protocol.licensezero.com/1.0.0-pre/bundle.json
//go:generate ./goifyschema Currency https://protocol.licensezero.com/1.0.0-pre/currency.json
//go:generate ./goifyschema Digest https://protocol.licensezero.com/1.0.0-pre/digest.json
//go:generate ./goifyschema Jurisdiction https://protocol.licensezero.com/1.0.0-pre/jurisdiction.json
//go:generate ./goifyschema Key https://protocol.licensezero.com/1.0.0-pre/key.json
//go:generate ./goifyschema Ledger https://protocol.licensezero.com/1.0.0-pre/ledger.json
//go:generate ./goifyschema Offer https://protocol.licensezero.com/1.0.0-pre/offer.json
//go:generate ./goifyschema Price https://protocol.licensezero.com/1.0.0-pre/price.json
//go:generate ./goifyschema Receipt https://protocol.licensezero.com/1.0.0-pre/receipt.json
//go:generate ./goifyschema Register https://protocol.licensezero.com/1.0.0-pre/register.json
//go:generate ./goifyschema Signature https://protocol.licensezero.com/1.0.0-pre/signature.json
//go:generate ./goifyschema Time https://protocol.licensezero.com/1.0.0-pre/time.json
//go:generate ./goifyschema URL https://protocol.licensezero.com/1.0.0-pre/url.json

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
