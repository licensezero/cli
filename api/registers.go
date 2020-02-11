package api

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
	"time"
)

// Register represents a register of broker signing keys.
type Register struct {
	Updated string               `json:"updated"`
	Keys    map[string]Timeframe `json:"keys"`
}

// Timeframe represents the period of time when a key is valid.
type Timeframe struct {
	From    *RegisterTime `json:"from"`
	Through *RegisterTime `json:"through,omitempty"`
}

type RegisterTime struct{ time.Time }

func (kt *RegisterTime) format() string {
	return kt.UTC().Format(time.RFC3339)
}

func (kt *RegisterTime) UnmarshalJSON(b []byte) (err error) {
	var asString string
	err = json.Unmarshal(b, &asString)
	if err != nil {
		return
	}
	t, err := time.Parse(time.RFC3339, asString)
	if err != nil {
		return
	}
	*kt = RegisterTime{t}
	return
}

func (kt *RegisterTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + kt.format() + "\""), nil
}

// ErrInvalidRegister indicates that a Register struct does not conform
// to the JSON schema for registers.
var ErrInvalidRegister = errors.New("invalid register")

var registerValidator *gojsonschema.Schema = nil

// Validate checks that a Register conforms to the
// JSON schema for register records.
func (register *Register) Validate() error {
	if registerValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Register),
		)
		if err != nil {
			panic(err)
		}
		registerValidator = schema
	}
	marshaled, err := json.Marshal(register)
	if err != nil {
		return err
	}
	dataLoader := gojsonschema.NewBytesLoader(marshaled)
	result, err := registerValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return ErrInvalidRegister
	}
	return nil
}
