package inventory

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseArtifact(t *testing.T) {
	server := "https://broker.licensezero.com"
	offerID := "a16ec460-0acb-4d5a-85e5-2787e61f084f"
	public := "Prosperity-3.0.0"
	data := fmt.Sprintf(`{"offers": [ { "server": "%v", "offerID": "%v", "public": "%v" } ] }`, server, offerID, public)

	var unstructured interface{}
	err := json.Unmarshal([]byte(data), &unstructured)
	if err != nil {
		t.Fatal(err)
	}

	artifact, err := mapToArtifact(unstructured)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifact.Offers) != 1 {
		t.Fatal("failed to parse one offer")
	}
	offer := artifact.Offers[0]
	if offer.Server != server {
		t.Error("failed to parse server")
	}
	if offer.OfferID != offerID {
		t.Error("failed to parse offer ID")
	}
	if offer.Public != public {
		t.Error("failed to public license")
	}
}
