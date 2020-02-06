package user

import (
	"encoding/hex"
	"encoding/json"
	"github.com/licensezero/helptest"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestReadReceipts(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	receipts := path.Join(directory, "receipts")
	err := os.MkdirAll(receipts, 0700)
	if err != nil {
		t.Fatal(err)
	}

	withVendor := path.Join(receipts, "withVendor.json")
	err = ioutil.WriteFile(
		withVendor,
		[]byte(`
{
  "key": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "signature": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "license": {
    "values": {
      "api": "https://api.licensezero.com",
      "offerID": "9aab7058-599a-43db-9449-5fc0971ecbfa",
      "effective": "2018-11-13T20:20:39Z",
      "expires": "2019-11-13T20:20:39Z",
      "orderID": "2c743a84-09ce-4549-9f0d-19d8f53462bb",
      "licensee": {
        "email": "licensee@example.com",
        "jurisdiction": "US-TX",
        "name": "Joe Licensee"
      },
      "licensor": {
        "email": "licensor@example.com",
        "jurisdiction": "US-CA",
        "name": "Jane Licensor",
        "licensorID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
      },
      "price": {
        "amount": 1000,
        "currency": "USD"
      },
      "vendor": {
        "email": "support@artlessdevices.com",
        "name": "Artless Devices LLC",
        "jurisdiction": "US-CA",
        "website": "https://artlessdevices.com"
      }
    },
    "form": "Test license form."
  }
}
			`),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	withoutVendor := path.Join(receipts, "withoutVendor.json")
	err = ioutil.WriteFile(
		withoutVendor,
		[]byte(`
{
  "key": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "signature": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "license": {
    "values": {
      "api": "https://api.licensezero.com",
      "offerID": "9aab7058-599a-43db-9449-5fc0971ecbfa",
      "effective": "2018-11-13T20:20:39Z",
      "orderID": "2c743a84-09ce-4549-9f0d-19d8f53462bb",
      "licensee": {
        "email": "licensee@example.com",
        "jurisdiction": "US-TX",
        "name": "Joe Licensee"
      },
      "licensor": {
        "email": "licensor@example.com",
        "jurisdiction": "US-CA",
        "name": "Jane Licensor",
        "licensorID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
      }
    },
    "form": "Test license form."
  }
}
			`),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	invalid := path.Join(receipts, "invalid.json")
	err = ioutil.WriteFile(invalid, []byte(`{}`), 0700)
	if err != nil {
		t.Fatal(err)
	}

	results, receiptErrors, readError := ReadReceipts(directory)
	if readError != nil {
		t.Fatal("read error")
	}

	if len(results) != 2 {
		t.Fatal("did not find receipt")
	}

	first := results[0]
	if first.API != "https://api.licensezero.com" {
		t.Error("failed to parse API")
	}
	if first.OrderID != "2c743a84-09ce-4549-9f0d-19d8f53462bb" {
		t.Error("failed to parse orderID")
	}
	if first.OfferID != "9aab7058-599a-43db-9449-5fc0971ecbfa" {
		t.Error("failed to parse orderID")
	}
	if first.Effective != "2018-11-13T20:20:39Z" {
		t.Error("failed to parse effective date")
	}
	if first.Expires != "2019-11-13T20:20:39Z" {
		t.Error("added expiration date")
	}

	if len(receiptErrors) != 1 {
		t.Error("missing invalid error")
	}
}

func TestValidateSignature(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	message := `{"form":"Test license form.","values":{"api":"https://api.licensezero.com","effective":"2018-11-13T20:20:39Z","licensee":{"email":"licensee@example.com","jurisdiction":"US-TX","name":"Joe"},"licensor":{"email":"licensor@example.com","jurisdiction":"US-CA","licensorID":"59e70a4d-ffee-4e9d-a526-7a9ff9161664","name":"Jane"},"offerID":"9aab7058-599a-43db-9449-5fc0971ecbfa","orderID":"2c743a84-09ce-4549-9f0d-19d8f53462bb"}}`
	signature := ed25519.Sign(privateKey, []byte(message))
	signatureHex := hex.EncodeToString(signature)
	publicKeyHex := hex.EncodeToString(publicKey)
	validJSON := "{" +
		"\"key\":" + "\"" + publicKeyHex + "\"" +
		",\"signature\":" + "\"" + signatureHex + "\"" +
		",\"license\":" + message +
		"}"
	var valid interface{}
	err := json.Unmarshal([]byte(validJSON), &valid)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateSignature(valid); err != nil {
		t.Error("invalidates invalid signature")
	}
	invalidJSON := strings.Replace(validJSON, publicKeyHex, strings.Repeat("a", 64), 1)
	var invalid interface{}
	err = json.Unmarshal([]byte(invalidJSON), &invalid)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateSignature(invalid); err == nil {
		t.Error("validates invalid signature")
	}
}
