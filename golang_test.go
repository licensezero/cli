package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFindGoDeps(t *testing.T) {
	withTestDir(t, func(directory string) {
		configDirectory := path.Join(directory, "config")
		err := os.MkdirAll(configDirectory, 0700)
		if err != nil {
			t.Fatal(err)
		}

		projectDirectory := path.Join(directory, "project")
		err = os.MkdirAll(projectDirectory, 0700)
		if err != nil {
			t.Fatal(err)
		}

		dep := "github.com/licensezero/gopackage"
		version := "0.0.1"
		err = ioutil.WriteFile(
			path.Join(projectDirectory, "go.mod"),
			[]byte(fmt.Sprintf(`
module licensezero.com/cli/golangtest
go 1.13
require %v v%v
`, dep, version)),
			0700,
		)
		if err != nil {
			t.Fatal(err)
		}

		err = ioutil.WriteFile(
			path.Join(projectDirectory, "main.go"),
			[]byte(fmt.Sprintf(`
package main

import "fmt"
import "%v"

func main () {
	fmt.Println(gopackage.Message)
}
			`, dep)),
			0700,
		)
		if err != nil {
			t.Fatal(err)
		}

		findings, err := findGoDeps(projectDirectory)
		if err != nil {
			t.Fatal("read error")
		}

		if len(findings) != 1 {
			t.Fatal("did not find one offer")
		}
		finding := findings[0]
		if finding.Name != "gopackage" {
			t.Error("did not set package name")
		}
		if finding.API != "https://api.commonform.com" {
			t.Error("did not set API")
		}
	})
}
