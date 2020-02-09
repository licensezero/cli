package main

import (
	"bytes"
	"errors"
	"flag"
	"github.com/licensezero/helptest"
	"io"
	"licensezero.com/licensezero/api"
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
	transport := helptest.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 404,
			Body:       helptest.NoopCloser{bytes.NewBufferString("")},
			Header:     make(http.Header),
		}
	})
	client := api.NewClient(transport)
	code = run(args, input, stdout, stderr, client)
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
