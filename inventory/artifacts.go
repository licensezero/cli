package inventory

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
	"licensezero.com/licensezero/schemas"
)

// Artifact encodes data about offers for an artifact.
type Artifact struct {
	Offers []ArtifactOffer `json:"offers" toml:"offers"`
}

// ArtifactOffer represents an offer relevant to an artifact.
type ArtifactOffer struct {
	API     string `json:"api" toml:"api"`
	OfferID string `json:"offerID" toml:"offerID"`
	Public  string `json:"public" toml:"public"`
}

func mapToArtifact(unstructred interface{}) (artifact Artifact, err error) {
	err = mapstructure.Decode(unstructred, artifact)
	return
}

var artifactValidator *gojsonschema.Schema = nil

// ErrInvalidArtifact indicates that the Artifact does not
// conform to the JSON schema for artifact metadata.
var ErrInvalidArtifact = errors.New("invalid artifact")

func (artifact *Artifact) Validate() error {
	if artifactValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Artifact),
		)
		if err != nil {
			panic(err)
		}
		artifactValidator = schema
	}
	marshaled, err := json.Marshal(artifact)
	if err != nil {
		return err
	}
	dataLoader := gojsonschema.NewBytesLoader(marshaled)
	result, err := artifactValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return ErrInvalidArtifact
	}
	return nil
}
