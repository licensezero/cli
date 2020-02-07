package user

import (
	"bytes"
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

func TestReceiptJSON(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	licenseJSON := `
{
  "form": "Test license form.",
  "values": {
    "api": "https://api.licensezero.com",
    "buyer": {
      "email": "buyer@example.com",
      "jurisdiction": "US-TX",
      "name": "Joe"
    },
    "effective": "2018-11-13T20:20:39Z",
    "offerID": "9aab7058-599a-43db-9449-5fc0971ecbfa",
    "orderID": "2c743a84-09ce-4549-9f0d-19d8f53462bb",
    "seller": {
      "email": "seller@example.com",
      "jurisdiction": "US-CA",
      "name": "Jane",
      "sellerID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
    }
  }
}
	`
	compactedLicenseJSON := bytes.NewBuffer([]byte{})
	err := json.Compact(compactedLicenseJSON, []byte(licenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	signature := ed25519.Sign(privateKey, compactedLicenseJSON.Bytes())
	signatureHex := hex.EncodeToString(signature)
	publicKeyHex := hex.EncodeToString(publicKey)
	receiptJSON := "{" +
		"\"key\":" + "\"" + publicKeyHex + "\"" +
		",\"license\":" + string(compactedLicenseJSON.Bytes()) +
		",\"signature\":" + "\"" + signatureHex + "\"" +
		"}"
	var receipt Receipt
	err = json.Unmarshal([]byte(receiptJSON), &receipt)
	if err != nil {
		t.Fatal(err)
	}
	reserialized, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, reserialized)
	if err != nil {
		t.Fatal(err)
	}
	if string(compacted.Bytes()) != receiptJSON {
		t.Error("serialization does not match compacted")
	}
	if err := ValidateReceipt(&receipt); err != nil {
		t.Error("invalidate valid receipt")
	}
	if err := ValidateSignature(&receipt); err != nil {
		t.Error("invalidates valid signature")
	}
}

func TestReadReceipts(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	receipts := path.Join(directory, "receipts")
	err := os.MkdirAll(receipts, 0700)
	if err != nil {
		t.Fatal(err)
	}

	withBroker := path.Join(receipts, "withBroker.json")
	err = ioutil.WriteFile(
		withBroker,
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
      "buyer": {
        "email": "buyer@example.com",
        "jurisdiction": "US-TX",
        "name": "Joe Buyer"
      },
      "seller": {
        "email": "seller@example.com",
        "jurisdiction": "US-CA",
        "name": "Jane Seller",
        "sellerID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
      },
      "price": {
        "amount": 1000,
        "currency": "USD"
      },
      "broker": {
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

	withoutBroker := path.Join(receipts, "withoutBroker.json")
	err = ioutil.WriteFile(
		withoutBroker,
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
      "buyer": {
        "email": "buyer@example.com",
        "jurisdiction": "US-TX",
        "name": "Joe Buyer"
      },
      "seller": {
        "email": "seller@example.com",
        "jurisdiction": "US-CA",
        "name": "Jane Seller",
        "sellerID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
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

	first := results[0].License.Values
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
	message := `{"form":"Test license form.","values":{"api":"https://api.licensezero.com","buyer":{"email":"buyer@example.com","jurisdiction":"US-TX","name":"Joe"},"effective":"2018-11-13T20:20:39Z","offerID":"9aab7058-599a-43db-9449-5fc0971ecbfa","orderID":"2c743a84-09ce-4549-9f0d-19d8f53462bb","seller":{"email":"seller@example.com","jurisdiction":"US-CA","name":"Jane","sellerID":"59e70a4d-ffee-4e9d-a526-7a9ff9161664"}}}`
	signature := ed25519.Sign(privateKey, []byte(message))
	signatureHex := hex.EncodeToString(signature)
	publicKeyHex := hex.EncodeToString(publicKey)
	validJSON := "{" +
		"\"key\":" + "\"" + publicKeyHex + "\"" +
		",\"signature\":" + "\"" + signatureHex + "\"" +
		",\"license\":" + message +
		"}"
	var valid Receipt
	err := json.Unmarshal([]byte(validJSON), &valid)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateSignature(&valid); err != nil {
		t.Error("invalidates invalid signature")
	}
	invalidJSON := strings.Replace(validJSON, publicKeyHex, strings.Repeat("a", 64), 1)
	var invalid Receipt
	err = json.Unmarshal([]byte(invalidJSON), &invalid)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateSignature(&invalid); err == nil {
		t.Error("validates invalid signature")
	}
}
