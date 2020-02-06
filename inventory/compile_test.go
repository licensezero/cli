package inventory

import (
	"bytes"
	"fmt"
	"github.com/licensezero/helptest"
	"io"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestCompileInventory(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	configDirectory := path.Join(directory, "config")
	err := os.MkdirAll(configDirectory, 0700)
	if err != nil {
		t.Fatal(err)
	}

	projectDirectory := path.Join(directory, "project")
	depDirectory := path.Join(projectDirectory, "dep")
	err = os.MkdirAll(depDirectory, 0700)
	if err != nil {
		t.Fatal(err)
	}

	vendorAPI := "https://api.licensezero.com"
	offerID := "186d34a9-c8f7-414c-91bc-a34b4553b91d"
	public := "Parity-7.0.0"
	offerURL := "http://example.com"
	err = ioutil.WriteFile(
		path.Join(depDirectory, "licensezero.json"),
		[]byte(fmt.Sprintf(`{"offers": [ { "api": "%v", "offerID": "%v", "public": "%v" } ] }`, vendorAPI, offerID, public)),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	licensorID := "902d42b2-77ee-4e67-aeb6-6dac3de8bb5a"
	licensorName := "Test Licensor"
	licensorJurisdiction := "US-CA"

	transport := testTransport(func(req *http.Request) *http.Response {
		url := req.URL.String()
		var json string
		if url == vendorAPI+"/offers/"+offerID {
			json = fmt.Sprintf(`
{
	"licensorID": "%v",
	"pricing": {
		"single": {
			"amount": 1000,
			"currency": "USD"
		}
	},
	"url": "%v"
}
			`, licensorID, offerURL)
			return &http.Response{
				StatusCode: 200,
				Body:       noopCloser{bytes.NewBufferString(json)},
				Header:     make(http.Header),
			}
		} else if url == vendorAPI+"/licensors/"+licensorID {
			json = fmt.Sprintf(
				`{ "name": "%v", "jurisdiction": "%v" }`,
				licensorName,
				licensorJurisdiction,
			)
			return &http.Response{
				StatusCode: 200,
				Body:       noopCloser{bytes.NewBufferString(json)},
				Header:     make(http.Header),
			}
		} else {
			return &http.Response{
				StatusCode: 404,
				Body:       noopCloser{bytes.NewBufferString("")},
				Header:     make(http.Header),
			}
		}
	})

	client := api.NewClient(transport)

	inventory, err := CompileInventory(
		configDirectory,
		projectDirectory,
		false,
		false,
		client,
	)
	if err != nil {
		t.Fatal("read error")
	}

	licensable := inventory.Licensable
	if len(licensable) != 1 {
		t.Fatal("did not find one licensable offer")
	}
	finding := licensable[0]
	if finding.API != vendorAPI {
		t.Error("did not read API")
	}
	if finding.OfferID != offerID {
		t.Error("did not read offer ID")
	}
	if finding.Public != public {
		t.Error("did not read public license")
	}
	if finding.Offer.LicensorID != licensorID {
		t.Error("did not read licensorID")
	}
	if finding.Offer.Pricing.Single.Amount != 1000 {
		t.Error("did not read offer single price amount")
	}
	if finding.Offer.Pricing.Single.Currency != "USD" {
		t.Error("did not read offer single price currency")
	}
	if finding.Offer.URL != offerURL {
		t.Error("did not read offer URL")
	}
	if finding.Licensor.Name != licensorName {
		t.Error("did not read licensor name")
	}
}

type testTransport func(req *http.Request) *http.Response

func (f testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type noopCloser struct {
	io.Reader
}

func (noopCloser) Close() error {
	return nil
}
