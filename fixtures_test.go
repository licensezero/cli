package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/licensezero/helptest"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

const offerID = "9aab7058-599a-43db-9449-5fc0971ecbfa"
const sellerID = "59e70a4d-ffee-4e9d-a526-7a9ff9161664"

var licenseJSON = fmt.Sprintf(`
{
  "form": "Test license form.",
  "values": {
    "api": "%v",
    "buyer": {
      "email": "buyer@example.com",
      "jurisdiction": "US-TX",
      "name": "Buyer"
    },
    "effective": "2018-11-13T20:20:39Z",
    "offerID": "%v",
    "orderID": "2c743a84-09ce-4549-9f0d-19d8f53462bb",
    "seller": {
      "email": "seller@example.com",
      "jurisdiction": "US-CA",
      "name": "Seller"
    },
    "sellerID": "%v"
  }
}`, "http://localhost:"+port, offerID, sellerID)

func writeGoodReceipt(t *testing.T) {
	t.Helper()
	if !*update {
		return
	}
	file := path.Join("testdata", "receipts", "good.json")
	os.MkdirAll(path.Join("testdata", "receipts"), 0755)
	ioutil.WriteFile(
		file,
		[]byte(generateGoodReceipt(t)),
		0644,
	)
	writeOffer(t)
}

func generateGoodReceipt(t *testing.T) string {
	t.Helper()
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	compactedLicenseJSON := generateReceiptLicense(t)
	signature := ed25519.Sign(privateKey, compactedLicenseJSON)
	signatureHex := hex.EncodeToString(signature)
	publicKeyHex := hex.EncodeToString(publicKey)
	return buildReceipt(
		publicKeyHex, signatureHex, compactedLicenseJSON,
	)
}

func writeOffer(t *testing.T) {
	t.Helper()
	if !*update {
		return
	}
	os.MkdirAll(path.Join("testdata", "offers"), 0755)
	ioutil.WriteFile(
		path.Join("testdata", "offers", offerID),
		[]byte(fmt.Sprintf(`{
"url": "http://example.com",
"sellerID": "%v",
"pricing": { "single": { "amount": 1000, "currency": "USD" } }
}`, sellerID)),
		0644,
	)
}

func generateReceiptLicense(t *testing.T) []byte {
	t.Helper()
	compactedLicenseJSON := bytes.NewBuffer([]byte{})
	err := json.Compact(compactedLicenseJSON, []byte(licenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	return compactedLicenseJSON.Bytes()
}

func buildReceipt(key, signature string, license []byte) string {
	return "{" +
		"\"key\":" + "\"" + key + "\"" +
		",\"license\":" + string(license) +
		",\"signature\":" + "\"" + signature + "\"" +
		"}"
}

func withTempConfig(t *testing.T) func() {
	t.Helper()
	directory, cleanup := helptest.TempDir(t, "licensezero")
	os.Setenv("LICENSEZERO_CONFIG", directory)
	return func() {
		cleanup()
	}
}

func writeBadReceipt(t *testing.T) {
	t.Helper()
	if !*update {
		return
	}
	os.MkdirAll(path.Join("testdata", "receipts"), 0755)
	ioutil.WriteFile(
		path.Join("testdata", "receipts", "bad.json"),
		[]byte(generateBadReceipt(t)),
		0644,
	)
}

func generateBadReceipt(t *testing.T) string {
	t.Helper()
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	compactedLicenseJSON := generateReceiptLicense(t)
	signature := ed25519.Sign(privateKey, compactedLicenseJSON)
	signatureHex := hex.EncodeToString(signature)
	publicKeyHex := hex.EncodeToString(publicKey)
	return buildReceipt(
		publicKeyHex,
		strings.Repeat("0", len(signatureHex)),
		compactedLicenseJSON,
	)
}

func writeGoodBundle(t *testing.T) {
	t.Helper()
	if !*update {
		return
	}
	json := "{" +
		"\"created\":\"2020-02-09T05:12:00Z\"," +
		"\"receipts\":[" + generateGoodReceipt(t) + "]" + "}"
	os.MkdirAll(path.Join("testdata", "bundles"), 0755)
	ioutil.WriteFile(
		path.Join("testdata", "bundles", "good.json"),
		[]byte(json),
		0644,
	)
	writeOffer(t)
}

func writeBadBundle(t *testing.T) {
	t.Helper()
	if !*update {
		return
	}
	json := "{" +
		"\"created\":\"2020-02-09T05:12:00Z\"," +
		"\"receipts\":[" + generateBadReceipt(t) + "]" + "}"
	os.MkdirAll(path.Join("testdata", "bundles"), 0755)
	ioutil.WriteFile(
		path.Join("testdata", "bundles", "bad.json"),
		[]byte(json),
		0644,
	)
}
