package api

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
)

// Keys represents a register of broker signing keys.
type Keys struct {
	Updated string               `json:"updated"`
	Keys    map[string]Timeframe `json:"keys"`
}

// Timeframe represents the period of time when a key is valid.
type Timeframe struct {
	From    string `json:"from"`
	Through string `json:"through"`
}

// ErrInvalidKeys indicates that a Keys struct does not conform
// to the JSON schema for keys records.
var ErrInvalidKeys = errors.New("invalid keys")

var keysValidator *gojsonschema.Schema = nil

// Validate checks that a Keys conforms to the
// JSON schema for keys records.
func (keys *Keys) Validate() error {
	if keysValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Keys),
		)
		if err != nil {
			panic(err)
		}
		keysValidator = schema
	}
	marshaled, err := json.Marshal(keys)
	if err != nil {
		return err
	}
	dataLoader := gojsonschema.NewBytesLoader(marshaled)
	result, err := keysValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return ErrInvalidKeys
	}
	return nil
}
