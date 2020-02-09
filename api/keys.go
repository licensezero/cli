package api

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
	"time"
)

// Keys represents a register of broker signing keys.
type Keys struct {
	Updated string               `json:"updated"`
	Keys    map[string]Timeframe `json:"keys"`
}

// Timeframe represents the period of time when a key is valid.
type Timeframe struct {
	From    *KeyTime `json:"from"`
	Through *KeyTime `json:"through,omitempty"`
}

type KeyTime struct{ time.Time }

func (kt *KeyTime) format() string {
	return kt.UTC().Format(time.RFC3339)
}

func (kt *KeyTime) UnmarshalJSON(b []byte) (err error) {
	var asString string
	err = json.Unmarshal(b, &asString)
	if err != nil {
		return
	}
	t, err := time.Parse(time.RFC3339, asString)
	if err != nil {
		return
	}
	*kt = KeyTime{t}
	return
}

func (kt *KeyTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + kt.format() + "\""), nil
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
