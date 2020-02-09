package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"net/http"
	"testing"
)

var update = flag.Bool("update", false, "update test fixtures")

type failingInputDevice struct{}

func (f *failingInputDevice) Confirm(string, io.StringWriter) (bool, error) {
	return false, errors.New("test input device")
}

func (f *failingInputDevice) SecretPrompt(string, io.StringWriter) (string, error) {
	return "", errors.New("test input device")
}

func runCommand(t *testing.T, args []string) (output string, errorOutput string, code int) {
	input := &failingInputDevice{}
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	code = run(args, input, stdout, stderr)
	output = string(stdout.Bytes())
	errorOutput = string(stderr.Bytes())
	return
}

const port = "8888"

func withTestDataServer(t *testing.T) func() {
	t.Helper()
	server := http.Server{
		Addr:    ":" + port,
		Handler: http.FileServer(http.Dir("testdata")),
	}
	go func() {
		server.ListenAndServe()
	}()
	return func() {
		server.Shutdown(nil)
	}
}
