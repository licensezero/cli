package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestCompileInventory(t *testing.T) {
	withTestDir(t, func(directory string) {
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

		inventory, err := compileInventory(configDirectory, projectDirectory, false, false)
		if err != nil {
			t.Fatal("read error")
		}

		licensable := inventory.Licensable
		if len(licensable) != 1 {
			t.Fatal("did not find one licensable offer")
		}
		finding := licensable[0]
		if finding.API != api {
			t.Error("did not read API")
		}
		if finding.OfferID != offerID {
			t.Error("did not read offer ID")
		}
		if finding.Public != public {
			t.Error("did not read public license")
		}
	})
}
