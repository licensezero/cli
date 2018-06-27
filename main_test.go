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

func TestWhoAmIWithoutIdentity(t *testing.T) {
	InTestDir(t, func() {
		command := exec.Command("./licensezero", "whoami")
		var stdout bytes.Buffer
		command.Stdout = &stdout
		err := command.Run()
		if err == nil {
			t.Error("Should fail")
		}
	})
}

func TestWhoAmIWithIdentity(t *testing.T) {
	InTestDir(t, func() {
		name := "John Doe"
		email := "test@example.com"
		jurisdiction := "US-CA"
		exec.Command("./licensezero", "identify", "--name", name, "--jurisdiction", jurisdiction, "--email", email, "--silent").Run()
		whoami := exec.Command("./licensezero", "whoami")
		var stdout bytes.Buffer
		whoami.Stdout = &stdout
		err := whoami.Run()
		if err != nil {
			t.Error(err)
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
	})
}

func TestImportWaiver(t *testing.T) {
	InTestDir(t, func() {
		name := "John Doe"
		email := "test@example.com"
		jurisdiction := "US-CA"
		exec.Command("./licensezero", "identify", "--name", name, "--jurisdiction", jurisdiction, "--email", email, "--silent").Run()
		importCommand := exec.Command("./licensezero", "import", "--file", "test/waiver.json")
		var stdout bytes.Buffer
		importCommand.Stdout = &stdout
		err := importCommand.Run()
		if err != nil {
			t.Error(err)
		}
		output := string(stdout.Bytes())
		if !strings.Contains(output, "Imported") {
			t.Error("does not say imported")
		}
	})
}

func TestImportNonexistentWaiver(t *testing.T) {
	InTestDir(t, func() {
		name := "John Doe"
		email := "test@example.com"
		jurisdiction := "US-CA"
		exec.Command("./licensezero", "identify", "--name", name, "--jurisdiction", jurisdiction, "--email", email, "--silent").Run()
		importCommand := exec.Command("./licensezero", "import", "--file", "test/nonexistent.json")
		var stdout bytes.Buffer
		importCommand.Stdout = &stdout
		err := importCommand.Run()
		if err == nil {
			t.Error("does not fail")
		}
	})
}

func TestImportLicense(t *testing.T) {
	InTestDir(t, func() {
		importCommand := exec.Command("./licensezero", "import", "--file", "test/license.json")
		var stdout bytes.Buffer
		importCommand.Stdout = &stdout
		err := importCommand.Run()
		if err != nil {
			t.Error(err)
		}
		output := string(stdout.Bytes())
		if !strings.Contains(output, "Imported") {
			t.Error("does not say imported")
		}
	})
}

func TestImportNonexistentLicense(t *testing.T) {
	InTestDir(t, func() {
		importCommand := exec.Command("./licensezero", "import", "--file", "test/nonexistent.json")
		var stdout bytes.Buffer
		importCommand.Stdout = &stdout
		err := importCommand.Run()
		if err == nil {
			t.Error("does not fail")
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
