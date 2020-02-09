package main

import (
	"strings"
	"testing"
)

func TestSanity(t *testing.T) {
	defer withTempConfig(t)
	output, _, code := runCommand(t, []string{})
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
	defer withTempConfig(t)
	output, _, _ := runCommand(t, []string{"identify", "--name", "John Doe", "--jurisdiction", "US-CA", "--email", "text@example.com"})
	if !strings.Contains(output, "Saved") {
		t.Error("Does not print \"Saved\"")
	}
}

func TestIdentifySilent(t *testing.T) {
	defer withTempConfig(t)
	output, _, _ := runCommand(t, []string{"identify", "--name", "John Doe", "--jurisdiction", "US-CA", "--email", "text@example.com", "--silent"})
	if output != "" {
		t.Error("No output")
	}
}

func TestWhoAmIWithoutIdentity(t *testing.T) {
	defer withTempConfig(t)
	_, _, code := runCommand(t, []string{"whoami"})
	if code == 0 {
		t.Error("exited with zero status")
	}
}

func TestWhoAmIWithIdentity(t *testing.T) {
	defer withTempConfig(t)
	name := "John Doe"
	email := "test@example.com"
	jurisdiction := "US-CA"
	runCommand(t, []string{"identify", "--name", name, "--jurisdiction", jurisdiction, "--email", email, "--silent"})
	output, _, _ := runCommand(t, []string{"whoami"})
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
	writeGoodReceipt(t)
	defer withTempConfig(t)
	output, _, _ := runCommand(t, []string{"import", "--file", "testdata/receipts/good.json"})
	if !strings.Contains(output, "Imported") {
		t.Error("does not say imported")
	}
}

func TestImportBadFile(t *testing.T) {
	writeBadReceipt(t)
	defer withTempConfig(t)
	_, errorOutput, _ := runCommand(t, []string{"import", "--file", "testdata/receipts/bad.json"})
	if !strings.Contains(errorOutput, "Invalid signature.") {
		t.Error("does not report invalid")
	}
}

func TestImportNonexistentFile(t *testing.T) {
	defer withTempConfig(t)
	_, _, code := runCommand(t, []string{"import", "--file", "testdata/receipts/nonexistent.json"})
	if code == 0 {
		t.Error("exited with zero status")
	}
}

func TestImportGoodBundle(t *testing.T) {
	writeGoodBundle(t)
	defer withTempConfig(t)
	defer withTestDataServer(t)
	output, errorOutput, _ := runCommand(t, []string{"import", "--bundle", "http://:" + port + "/bundles/good.json"})
	if errorOutput != "" {
		t.Error("error output")
	}
	if !strings.Contains(output, "Imported 1 licenses.") {
		t.Error("does not report imported")
	}
}

func TestImportBundleBadSignature(t *testing.T) {
	writeBadBundle(t)
	defer withTempConfig(t)
	defer withTestDataServer(t)
	output, errorOutput, code := runCommand(t, []string{"import", "--bundle", "http://:" + port + "/receipts/bad.json"})
	if code != 0 {
		t.Error("exited non-zero")
	}
	if !strings.Contains(output, "Imported 0 licenses.") {
		t.Error("does not report imported")
	}
	if !strings.Contains(errorOutput, "Invalid license signature") {
		t.Error("does not report invalid signature")
	}
}
