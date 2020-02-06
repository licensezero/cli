package inventory

import (
	"fmt"
	"github.com/licensezero/helptest"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFindLicenseZeroFiles(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	depDirectory := path.Join(directory, "dep")
	err := os.MkdirAll(depDirectory, 0700)
	if err != nil {
		t.Fatal(err)
	}

	api := "https://api.licensezero.com"
	offerID := "186d34a9-c8f7-414c-91bc-a34b4553b91d"
	public := "Parity-7.0.0"
	err = ioutil.WriteFile(
		path.Join(depDirectory, "licensezero.json"),
		[]byte(fmt.Sprintf(`{"offers": [ { "api": "%v", "offerID": "%v", "public": "%v" } ] }`, api, offerID, public)),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	findings, err := findLicenseZeroFiles(directory)
	if err != nil {
		t.Fatal("read error")
	}

	if len(findings) != 1 {
		t.Fatal("failed to find one offer")
	}
	finding := findings[0]
	if finding.API != api {
		t.Error("failed to parse API")
	}
	if finding.OfferID != offerID {
		t.Error("failed to parse offer ID")
	}
	if finding.Public != public {
		t.Error("failed to parse public license")
	}
}

func TestFindLicenseZeroFilesInNPMPackage(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	depDirectory := path.Join(directory, "node_modules", "dep")
	err := os.MkdirAll(depDirectory, 0700)
	if err != nil {
		t.Fatal(err)
	}

	name := "dep"
	version := "0.0.0"
	err = ioutil.WriteFile(
		path.Join(depDirectory, "package.json"),
		[]byte(fmt.Sprintf(`{"name": "%v", "version": "%v" }`, name, version)),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	api := "https://api.licensezero.com"
	offerID := "186d34a9-c8f7-414c-91bc-a34b4553b91d"
	public := "Parity-7.0.0"
	err = ioutil.WriteFile(
		path.Join(depDirectory, "licensezero.json"),
		[]byte(fmt.Sprintf(`{"offers": [ { "api": "%v", "offerID": "%v", "public": "%v" } ] }`, api, offerID, public)),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	findings, err := findLicenseZeroFiles(directory)
	if err != nil {
		t.Fatal("read error")
	}

	if len(findings) != 1 {
		t.Fatal("failed to find one offer")
	}
	finding := findings[0]
	if finding.Type != "npm" {
		t.Error("failed to mark as npm type")
	}
	if finding.Name != name {
		t.Error("failed to set name")
	}
	if finding.Version != version {
		t.Error("failed to set version")
	}
	if finding.Path != depDirectory {
		t.Error("failed to set path")
	}
	if finding.API != api {
		t.Error("failed to parse API")
	}
	if finding.OfferID != offerID {
		t.Error("failed to parse offer ID")
	}
	if finding.Public != public {
		t.Error("failed to parse public license")
	}
}
