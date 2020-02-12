package inventory

import (
	"github.com/licensezero/helptest"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestCargoCrates(t *testing.T) {
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

	err = ioutil.WriteFile(
		path.Join(projectDirectory, "Cargo.toml"),
		[]byte(`
[package]
name = "test"
version = "0.0.0"
publish = false

[dependencies]
license-zero-test-crate = { git = "https://github.com/licensezero/license-zero-test-crate#master" }
			`),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(
		path.Join(srcDirectory, "main.rs"),
		[]byte(`
extern crate license_zero_test_crate as test;

fn main () {
	println!("{}", test::MESSAGE);
}
			`),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	build := exec.Command("cargo", "build")
	build.Dir = projectDirectory
	err = build.Run()
	if err != nil {
		t.Fatal(err)
	}

	findings, err := findCargoCrates(projectDirectory)
	if err != nil {
		t.Fatal("read error")
	}

	if len(findings) != 1 {
		t.Fatal("did not find one offer")
	}
	finding := findings[0]
	if finding.Type != "cargo" {
		t.Error("did not set Cargo type")
	}
	if finding.Server != "https://broker.licensezero.com" {
		t.Log(finding)
		t.Error("did not set server")
	}
}
