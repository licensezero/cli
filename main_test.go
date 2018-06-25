package main

import "bytes"
import "io/ioutil"
import "os"
import "os/exec"
import "strings"
import "testing"

func TestSanity(t *testing.T) {
	command := exec.Command("./licensezero")
	var stdout bytes.Buffer
	command.Stdout = &stdout
	err := command.Run()
	if err != nil {
		t.Error(err)
	}
	output := string(stdout.Bytes())
	if !strings.Contains(output, "Subcommands:") {
		t.Error("does not list subcommands")
	}
	if !strings.Contains(output, "License Zero") {
		t.Error("does not mention License Zero")
	}
}

func TestIdentify(t *testing.T) {
	InTestDir(t, func() {
		command := exec.Command("./licensezero", "identify", "--name", "John Doe", "--jurisdiction", "US-CA", "--email", "text@example.com")
		var stdout bytes.Buffer
		command.Stdout = &stdout
		err := command.Run()
		if err != nil {
			t.Error(err)
		}
		output := string(stdout.Bytes())
		if !strings.Contains(output, "Saved") {
			t.Error("Does not print \"Saved\"")
		}
	})
}

func TestIdentifySilent(t *testing.T) {
	InTestDir(t, func() {
		command := exec.Command("./licensezero", "identify", "--name", "John Doe", "--jurisdiction", "US-CA", "--email", "text@example.com", "--silent")
		var stdout bytes.Buffer
		command.Stdout = &stdout
		err := command.Run()
		if err != nil {
			t.Error(err)
		}
		output := string(stdout.Bytes())
		if output != "" {
			t.Error("No output")
		}
	})
}

func InTestDir(t *testing.T, script func()) {
	directory, err := ioutil.TempDir("/tmp", "licensezero-test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(directory)
	os.Setenv("LICENSEZERO_CONFIG", directory)
	script()
}
