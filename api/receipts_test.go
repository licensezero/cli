package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/ed25519"
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
      "name": "Jane"
    },
    "sellerID": "59e70a4d-ffee-4e9d-a526-7a9ff9161664"
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
	t.Log(string(compacted.Bytes()))
	t.Log(receiptJSON)
	if string(compacted.Bytes()) != receiptJSON {
		t.Error("serialization does not match compacted")
	}
	if err := receipt.Validate(); err != nil {
		t.Error("invalidate valid receipt")
	}
	if err := receipt.VerifySignature(); err != nil {
		t.Error("invalidates valid signature")
	}
}

func TestValidateSignature(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	message := `{"form":"Test license form.","values":{"api":"https://api.licensezero.com","buyer":{"email":"buyer@example.com","jurisdiction":"US-TX","name":"Joe"},"effective":"2018-11-13T20:20:39Z","offerID":"9aab7058-599a-43db-9449-5fc0971ecbfa","orderID":"2c743a84-09ce-4549-9f0d-19d8f53462bb","seller":{"email":"seller@example.com","jurisdiction":"US-CA","name":"Jane"},"sellerID":"59e70a4d-ffee-4e9d-a526-7a9ff9161664"}}`
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
	if err := valid.VerifySignature(); err != nil {
		t.Error("invalidates invalid signature")
	}
	invalidJSON := strings.Replace(validJSON, publicKeyHex, strings.Repeat("a", 64), 1)
	var invalid Receipt
	err = json.Unmarshal([]byte(invalidJSON), &invalid)
	if err != nil {
		t.Fatal(err)
	}
	if err := invalid.VerifySignature(); err == nil {
		t.Error("validates invalid signature")
	}
}
