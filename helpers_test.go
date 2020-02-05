package cli

import (
	"io/ioutil"
	"os"
	"testing"
)

func withTestDir(t *testing.T, script func(string)) {
	t.Helper()
	directory, err := ioutil.TempDir("/tmp", "licensezero-test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(directory)
	script(directory)
}
