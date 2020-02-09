package user

import (
	"github.com/licensezero/helptest"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestReadIdentity(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	email := "test@example.com"
	jurisdiction := "US-CA"
	name := "D Tester"
	err := ioutil.WriteFile(
		path.Join(directory, "identity.json"),
		[]byte("{\"email\": \""+email+"\", \"jurisdiction\": \""+jurisdiction+"\", \"name\": \""+name+"\"}"),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("LICENSEZERO_CONFIG", directory)

	result, err := ReadIdentity()
	if err != nil {
		t.Fatal("read error")
	}

	if result.Name != name {
		t.Error("did not read name")
	}
	if result.Jurisdiction != jurisdiction {
		t.Error("did not read jurisdiction")
	}
	if result.EMail != email {
		t.Error("did not read e-mail")
	}
}
