package main

import (
	"bytes"
	"errors"
	"flag"
	"github.com/licensezero/helptest"
	"io"
	"net/http"
	"os"
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
			Body:       helptest.NoopCloser{Reader: bytes.NewBufferString("")},
			Header:     make(http.Header),
		}
	})
	client := http.Client{Transport: transport}
	code = run(args, input, stdout, stderr, &client)
	output = string(stdout.Bytes())
	errorOutput = string(stderr.Bytes())
	return
}

func mockClient(t *testing.T, responses map[string]string) *http.Client {
	transport := helptest.RoundTripFunc(func(req *http.Request) *http.Response {
		url := req.URL.String()
		json, ok := responses[url]
		if !ok {
			return &http.Response{
				StatusCode: 404,
				Body:       helptest.NoopCloser{Reader: bytes.NewBufferString("")},
				Header:     make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: 200,
			Body:       helptest.NoopCloser{Reader: bytes.NewBufferString(json)},
			Header:     make(http.Header),
		}
	})
	return &http.Client{Transport: transport}
}

func withTempConfig(t *testing.T) func() {
	t.Helper()
	directory, cleanup := helptest.TempDir(t, "licensezero")
	os.Setenv("LICENSEZERO_CONFIG", directory)
	return func() {
		cleanup()
	}
}
