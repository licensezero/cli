package cli

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseOffer(t *testing.T) {
	url := "http://example.com/project"
	licensorID := "6005a0cd-2481-468c-a53b-6e844930e413"
	var amount uint = 500
	currency := "USD"
	data := fmt.Sprintf(`{"url": "%v", "licensorID": "%v", "pricing": { "single": { "amount": %v, "currency": "%v" } } }`, url, licensorID, amount, currency)

	var unstructured interface{}
	err := json.Unmarshal([]byte(data), &unstructured)
	if err != nil {
		t.Fatal(err)
	}

	offer, err := parseOffer(unstructured)
	if err != nil {
		t.Fatal(err)
	}
	if offer.URL != url {
		t.Error("failed to parse URL")
	}
	if offer.LicensorID != licensorID {
		t.Error("failed to parse licensor ID")
	}
	if offer.Pricing.Single.Amount != amount {
		t.Error("failed to parse amount")
	}
	if offer.Pricing.Single.Currency != currency {
		t.Error("failed to parse currency")
	}
}
