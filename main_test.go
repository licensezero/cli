package main

import "bytes"
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
