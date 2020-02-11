package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestSanity(t *testing.T) {
	defer withTempConfig(t)
	output, _, code := runCommand(t, []string{})
	if code != 0 {
		t.Error("exited with non-zero status")
	}
	if !strings.Contains(output, "Subcommands:") {
		t.Error("does not list subcommands")
	}
	if !strings.Contains(output, "License Zero") {
		t.Error("does not mention License Zero")
	}
}

func TestIdentify(t *testing.T) {
	defer withTempConfig(t)
	output, _, _ := runCommand(t, []string{"identify", "--name", "John Doe", "--jurisdiction", "US-CA", "--email", "text@example.com"})
	if !strings.Contains(output, "Saved") {
		t.Error("Does not print \"Saved\"")
	}
}

func TestIdentifySilent(t *testing.T) {
	defer withTempConfig(t)
	output, _, _ := runCommand(t, []string{"identify", "--name", "John Doe", "--jurisdiction", "US-CA", "--email", "text@example.com", "--silent"})
	if output != "" {
		t.Error("No output")
	}
}

func TestWhoAmIWithoutIdentity(t *testing.T) {
	defer withTempConfig(t)
	_, _, code := runCommand(t, []string{"whoami"})
	if code == 0 {
		t.Error("exited with zero status")
	}
}

func TestWhoAmIWithIdentity(t *testing.T) {
	defer withTempConfig(t)
	name := "John Doe"
	email := "test@example.com"
	jurisdiction := "US-CA"
	runCommand(t, []string{"identify", "--name", name, "--jurisdiction", jurisdiction, "--email", email, "--silent"})
	output, _, _ := runCommand(t, []string{"whoami"})
	if !strings.Contains(output, name) {
		t.Error("does not list name")
	}
	if !strings.Contains(output, email) {
		t.Error("does not list e-mail")
	}
	if !strings.Contains(output, jurisdiction) {
		t.Error("does not list jurisdiction")
	}
}

func TestImportGoodFile(t *testing.T) {
	writeGoodReceipt(t)
	defer withTempConfig(t)
	output, _, _ := runCommand(t, []string{"import", "--file", "testdata/good-receipt.json"})
	if !strings.Contains(output, "Imported") {
		t.Error("does not say imported")
	}
}

func writeGoodReceipt(t *testing.T) {
	t.Helper()
	if !*update {
		return
	}
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	offerID := uuid.New().String()
	orderID := uuid.New().String()
	sellerID := uuid.New().String()
	effective := time.Now().UTC().Format(time.RFC3339)
	license := api.License{
		Form: "test form",
		Values: api.Values{
			API: "https://broker.licensezero.com",
			Buyer: &api.Buyer{
				EMail:        "buyer@example.com",
				Jurisdiction: "US-TX",
				Name:         "Buyer",
			},
			Effective: effective,
			OfferID:   offerID,
			OrderID:   orderID,
			SellerID:  sellerID,
			Seller: &api.Seller{
				EMail:        "seller@example.com",
				Jurisdiction: "US-CA",
				Name:         "Seller",
			},
		},
	}
	licenseJSON, err := json.Marshal(license)
	if err != nil {
		t.Fatal(err)
	}
	compactedLicenseJSON := bytes.NewBuffer([]byte{})
	err = json.Compact(compactedLicenseJSON, []byte(licenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	publicKeyHex := hex.EncodeToString(publicKey)
	signature := ed25519.Sign(privateKey, compactedLicenseJSON.Bytes())
	signatureHex := hex.EncodeToString(signature)
	receipt := api.Receipt{
		KeyHex:       publicKeyHex,
		License:      license,
		SignatureHex: signatureHex,
	}
	receiptJSON, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll("testdata", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(
		path.Join("testdata", "good-receipt.json"),
		receiptJSON,
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestImportBadFile(t *testing.T) {
	writeBadReceipt(t)
	defer withTempConfig(t)
	_, errorOutput, _ := runCommand(t, []string{"import", "--file", "testdata/bad-receipt.json"})
	if !strings.Contains(errorOutput, "Invalid signature.") {
		t.Error("does not report invalid")
	}
}

func writeBadReceipt(t *testing.T) {
	t.Helper()
	if !*update {
		return
	}
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	offerID := uuid.New().String()
	orderID := uuid.New().String()
	sellerID := uuid.New().String()
	effective := time.Now().UTC().Format(time.RFC3339)
	license := api.License{
		Form: "test form",
		Values: api.Values{
			API: "https://broker.licensezero.com",
			Buyer: &api.Buyer{
				EMail:        "buyer@example.com",
				Jurisdiction: "US-TX",
				Name:         "Buyer",
			},
			Effective: effective,
			OfferID:   offerID,
			OrderID:   orderID,
			SellerID:  sellerID,
			Seller: &api.Seller{
				EMail:        "seller@example.com",
				Jurisdiction: "US-CA",
				Name:         "Seller",
			},
		},
	}
	licenseJSON, err := json.Marshal(license)
	if err != nil {
		t.Fatal(err)
	}
	compactedLicenseJSON := bytes.NewBuffer([]byte{})
	err = json.Compact(compactedLicenseJSON, []byte(licenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	publicKeyHex := hex.EncodeToString(publicKey)
	signature := ed25519.Sign(privateKey, compactedLicenseJSON.Bytes())
	signatureHex := hex.EncodeToString(signature)
	invalidSignatureHex := strings.Repeat("0", len(signatureHex))
	receipt := api.Receipt{
		KeyHex:       publicKeyHex,
		License:      license,
		SignatureHex: invalidSignatureHex,
	}
	receiptJSON, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll("testdata", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(
		path.Join("testdata", "bad-receipt.json"),
		receiptJSON,
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestImportNonexistentFile(t *testing.T) {
	defer withTempConfig(t)
	_, _, code := runCommand(t, []string{"import", "--file", "testdata/receipts/nonexistent.json"})
	if code == 0 {
		t.Error("exited with zero status")
	}
}

func TestImportGoodBundle(t *testing.T) {
	defer withTempConfig(t)

	// Generate receipt and other records.
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	brokerAPI := "https://broker.licensezero.com"
	offerID := uuid.New().String()
	orderID := uuid.New().String()
	sellerID := uuid.New().String()
	bundleID := uuid.New().String()
	bundleURL := brokerAPI + "/receipts/" + bundleID
	effective := time.Now().UTC().Format(time.RFC3339)
	license := api.License{
		Form: "test form",
		Values: api.Values{
			API: brokerAPI,
			Buyer: &api.Buyer{
				EMail:        "buyer@example.com",
				Jurisdiction: "US-TX",
				Name:         "Buyer",
			},
			Effective: effective,
			OfferID:   offerID,
			OrderID:   orderID,
			SellerID:  sellerID,
			Seller: &api.Seller{
				EMail:        "seller@example.com",
				Jurisdiction: "US-CA",
				Name:         "Seller",
			},
		},
	}
	licenseJSON, err := json.Marshal(license)
	if err != nil {
		t.Fatal(err)
	}
	compactedLicenseJSON := bytes.NewBuffer([]byte{})
	err = json.Compact(compactedLicenseJSON, []byte(licenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	signature := ed25519.Sign(privateKey, compactedLicenseJSON.Bytes())
	signatureHex := hex.EncodeToString(signature)
	publicKeyHex := hex.EncodeToString(publicKey)
	receipt := api.Receipt{
		KeyHex:       publicKeyHex,
		License:      license,
		SignatureHex: signatureHex,
	}
	receiptJSON, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}

	// Run import subcommand.
	input := &failingInputDevice{}
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	arguments := []string{"import", "--bundle", bundleURL}
	client := mockClient(t, map[string]string{
		bundleURL: fmt.Sprintf(
			`{ "created": "%v", "receipts": [%v] }`,
			effective, string(receiptJSON),
		),
	})
	code := run(arguments, input, stdout, stderr, client)
	output := string(stdout.Bytes())
	if code != 0 {
		t.Error("exited non-zero")
	}
	if strings.Contains(output, "Imported 1 receipt.") {
		t.Error("does not report imported")
	}
}

func TestImportBundleBadSignature(t *testing.T) {
}
