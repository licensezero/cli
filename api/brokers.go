package api

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
)

// Broker represents a party selling a license on behalf of a Seller.
type Broker struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
	Website      string `json:"website"`
}

// ErrInvalidBroker indicates that a Broker struct does not conform
// to the JSON schema for brokers.
var ErrInvalidBroker = errors.New("invalid broker")

var brokerValidator *gojsonschema.Schema = nil

// Validate checks that a Broker conforms to the
// JSON schema for broker records.
func (broker *Broker) Validate() error {
	if brokerValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Broker),
		)
		if err != nil {
			panic(err)
		}
		brokerValidator = schema
	}
	marshaled, err := json.Marshal(broker)
	if err != nil {
		return err
	}
	dataLoader := gojsonschema.NewBytesLoader(marshaled)
	result, err := brokerValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return ErrInvalidBroker
	}
	return nil
}
