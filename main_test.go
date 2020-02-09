package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/licensezero/helptest"
	"golang.org/x/crypto/ed25519"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update test fixtures")

type failingInputDevice struct{}

func (f *failingInputDevice) Confirm(string, io.StringWriter) (bool, error) {
	return false, errors.New("test input device")
}

func (f *failingInputDevice) SecretPrompt(string, io.StringWriter) (string, error) {
	return "", errors.New("test input device")
}

func TestSanity(t *testing.T) {
	input := &failingInputDevice{}
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	code := run(
		[]string{},
		input,
		stdout,
		stderr,
	)
	output := string(stdout.Bytes())
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
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	command := exec.Command("./licensezero", "identify", "--name", "John Doe", "--jurisdiction", "US-CA", "--email", "text@example.com")
	command.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	var stdout bytes.Buffer
	command.Stdout = &stdout
	err := command.Run()
	if err != nil {
		t.Fatal(err)
	}
	output := string(stdout.Bytes())
	if !strings.Contains(output, "Saved") {
		t.Error("Does not print \"Saved\"")
	}
}

func TestIdentifySilent(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	command := exec.Command("./licensezero", "identify", "--name", "John Doe", "--jurisdiction", "US-CA", "--email", "text@example.com", "--silent")
	command.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	var stdout bytes.Buffer
	command.Stdout = &stdout
	err := command.Run()
	if err != nil {
		t.Fatal(err)
	}
	output := string(stdout.Bytes())
	if output != "" {
		t.Error("No output")
	}
}

func TestWhoAmIWithoutIdentity(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	command := exec.Command("./licensezero", "whoami")
	command.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	var stdout bytes.Buffer
	command.Stdout = &stdout
	err := command.Run()
	if err == nil {
		t.Error("Should fail")
	}
}

func TestWhoAmIWithIdentity(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	name := "John Doe"
	email := "test@example.com"
	jurisdiction := "US-CA"
	identify := exec.Command("./licensezero", "identify", "--name", name, "--jurisdiction", jurisdiction, "--email", email, "--silent")
	identify.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	err := identify.Run()
	if err != nil {
		t.Fatal(err)
	}
	whoami := exec.Command("./licensezero", "whoami")
	whoami.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	var stdout bytes.Buffer
	whoami.Stdout = &stdout
	err = whoami.Run()
	if err != nil {
		t.Fatal(err)
	}
	output := string(stdout.Bytes())
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
	if *update {
		writeGoodReceipt(t)
	}
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	importCommand := exec.Command("./licensezero", "import", "--file", "testdata/receipts/good.json")
	importCommand.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	var stdout bytes.Buffer
	importCommand.Stdout = &stdout
	err := importCommand.Run()
	if err != nil {
		t.Fatal(err)
	}
	output := string(stdout.Bytes())
	if !strings.Contains(output, "Imported") {
		t.Error("does not say imported")
	}
}

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

func TestImportBadFile(t *testing.T) {
	if *update {
		writeBadReceipt(t)
	}
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	command := exec.Command("./licensezero", "import", "--file", "testdata/receipts/bad.json")
	command.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err == nil {
		t.Error("does not fail")
	}
	if !strings.Contains(string(stderr.Bytes()), "Invalid signature.") {
		t.Error("does not report invalid")
	}
}

func writeBadReceipt(t *testing.T) {
	t.Helper()
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

func TestImportNonexistentFile(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	importCommand := exec.Command("./licensezero", "import", "--file", "testdata/receipts/nonexistent.json")
	importCommand.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	var stdout bytes.Buffer
	importCommand.Stdout = &stdout
	err := importCommand.Run()
	if err == nil {
		t.Error("does not fail")
	}
}

func TestImportGoodBundle(t *testing.T) {
	if *update {
		writeGoodBundle(t)
	}
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	defer withTestDataServer(t)
	command := exec.Command("./licensezero", "import", "--bundle", "http://:"+port+"/bundles/good.json")
	command.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		t.Error("does not fail")
	}
	if string(stderr.Bytes()) != "" {
		t.Error("error output")
	}
	if !strings.Contains(string(stdout.Bytes()), "Imported 1 licenses.") {
		t.Error("does not report imported")
	}
}

func writeGoodBundle(t *testing.T) {
	t.Helper()
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

func TestImportBundleBadSignature(t *testing.T) {
	if *update {
		writeBadBundle(t)
	}
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	defer withTestDataServer(t)
	command := exec.Command("./licensezero", "import", "--bundle", "http://:"+port+"/receipts/bad.json")
	command.Env = []string{"LICENSEZERO_CONFIG=" + directory}
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		t.Error("does not fail")
	}
	if !strings.Contains(string(stdout.Bytes()), "Imported 0 licenses.") {
		t.Error("does not report imported")
	}
	if !strings.Contains(string(stderr.Bytes()), "Invalid license signature") {
		t.Error("does not report invalid signature")
	}
}

func writeBadBundle(t *testing.T) {
	t.Helper()
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

func Identify() {
	name := "John Doe"
	email := "test@example.com"
	jurisdiction := "US-CA"
	exec.Command("./licensezero", "identify", "--name", name, "--jurisdiction", jurisdiction, "--email", email, "--silent").Run()
}

const port = "8888"

func withTestDataServer(t *testing.T) func() {
	t.Helper()
	server := http.Server{
		Addr:    ":" + port,
		Handler: http.FileServer(http.Dir("testdata")),
	}
	go func() {
		server.ListenAndServe()
	}()
	return func() {
		server.Shutdown(nil)
	}
}
