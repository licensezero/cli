package inventory

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/licensezero/helptest"
)

func TestCompile(t *testing.T) {
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

	brokerServer := "https://broker.licensezero.com"
	offerID := "186d34a9-c8f7-414c-91bc-a34b4553b91d"
	public := "Parity-7.0.0"
	offerURL := "https://example.com"
	err = ioutil.WriteFile(
		path.Join(depDirectory, "licensezero.json"),
		[]byte(fmt.Sprintf(`{"offers": [ { "server": "%v", "offerID": "%v", "public": "%v" } ] }`, brokerServer, offerID, public)),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	sellerID := "902d42b2-77ee-4e67-aeb6-6dac3de8bb5a"
	sellerEMail := "seller@example.com"
	sellerName := "Test Seller"
	sellerJurisdiction := "US-CA"
	sellerURL := "https://example.com/~seller"

	transport := helptest.RoundTripFunc(func(req *http.Request) *http.Response {
		url := req.URL.String()
		var json string
		if url == brokerServer+"/offers/"+offerID {
			json = fmt.Sprintf(`
{
	"sellerID": "%v",
	"pricing": {
		"single": {
			"amount": 1000,
			"currency": "USD"
		}
	},
	"url": "%v"
}
			`, sellerID, offerURL)
			return &http.Response{
				StatusCode: 200,
				Body:       helptest.NoopCloser{Reader: bytes.NewBufferString(json)},
				Header:     make(http.Header),
			}
		} else if url == brokerServer+"/sellers/"+sellerID {
			json = fmt.Sprintf(
				`{ "email": "%v", "name": "%v", "jurisdiction": "%v", "url": "%v" }`,
				sellerEMail,
				sellerName,
				sellerJurisdiction,
				sellerURL,
			)
			return &http.Response{
				StatusCode: 200,
				Body:       helptest.NoopCloser{Reader: bytes.NewBufferString(json)},
				Header:     make(http.Header),
			}
		} else {
			return &http.Response{
				StatusCode: 404,
				Body:       helptest.NoopCloser{Reader: bytes.NewBufferString("")},
				Header:     make(http.Header),
			}
		}
	})

	client := http.Client{Transport: transport}

	inventory, err := Compile(
		configDirectory,
		projectDirectory,
		false,
		false,
		&client,
	)
	if err != nil {
		t.Fatal("read error")
	}

	licensable := inventory.Licensable
	if len(licensable) != 1 {
		t.Fatal("did not find one licensable offer")
	}
	finding := licensable[0]
	if finding.Server != brokerServer {
		t.Error("did not read Server")
	}
	if finding.OfferID != offerID {
		t.Error("did not read offer ID")
	}
	if finding.Public != public {
		t.Error("did not read public license")
	}
	if finding.Offer.SellerID != sellerID {
		t.Error("did not read sellerID")
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
	if finding.Seller.Name != sellerName {
		t.Error("did not read seller name")
	}
	if finding.Seller.URL != sellerURL {
		t.Error("did not read seller URL")
	}
}
