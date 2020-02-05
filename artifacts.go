package cli

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/xeipuuv/gojsonschema"
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

func parseArtifact(unstructured interface{}) (a *Artifact, err error) {
	err = validateArtifact(unstructured)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(unstructured, &a)
	return
}

var artifactValidator *gojsonschema.Schema = nil

func validateArtifact(unstructured interface{}) error {
	if artifactValidator == nil {
		schema, err := schemaLoader().Compile(
			gojsonschema.NewStringLoader(artifactSchema),
		)
		if err != nil {
			panic(err)
		}
		artifactValidator = schema
	}
	dataLoader := gojsonschema.NewGoLoader(unstructured)
	result, err := artifactValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return errors.New("invalid artifact")
	}
	return nil
}
