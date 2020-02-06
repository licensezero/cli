package inventory

import (
	"fmt"
	"github.com/licensezero/helptest"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestNPM(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	configDirectory := path.Join(directory, "config")
	err := os.MkdirAll(configDirectory, 0700)
	if err != nil {
		t.Fatal(err)
	}

	projectDirectory := path.Join(directory, "project")
	srcDirectory := path.Join(projectDirectory, "src")
	err = os.MkdirAll(srcDirectory, 0700)
	if err != nil {
		t.Fatal(err)
	}

	name := "example-licensezero-package"
	err = ioutil.WriteFile(
		path.Join(projectDirectory, "package.json"),
		[]byte(fmt.Sprintf(`
{
	"private": true,
	"dependencies": {
		"%v": "git+https://github.com/licensezero/%v"
	}
}
			`, name, name)),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	install := exec.Command("npm", "install")
	install.Dir = projectDirectory
	err = install.Run()
	if err != nil {
		t.Fatal(err)
	}

	findings, err := findLicenseZeroFiles(projectDirectory)
	if err != nil {
		t.Fatal("read error")
	}

	if len(findings) != 1 {
		t.Fatal("did not find one offer")
	}
	finding := findings[0]
	if finding.Name != name {
		t.Error("did not set name")
	}
	if finding.Type != "npm" {
		t.Error("did not set npm type")
	}
	if finding.API != "https://api.licensezero.com" {
		t.Error("did not set API")
	}
	if finding.Path != path.Join(projectDirectory, "node_modules", name) {
		t.Error("did not set path")
	}
}
