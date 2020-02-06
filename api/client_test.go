package api

import (
	"bytes"
	"github.com/licensezero/helptest"
	"net/http"
	"testing"
)

func TestInvalidResponses(t *testing.T) {
	vendorAPI := "https://licensezero.com"
	offerID := "186d34a9-c8f7-414c-91bc-a34b4553b91d"
	licensorID := "2d7aa023-8eca-4d15-ad43-6b15768a5293"
	transport := helptest.RoundTripFunc(func(req *http.Request) *http.Response {
		url := req.URL.String()
		if url == vendorAPI+"/offers/"+offerID ||
			url == vendorAPI+"/licensors/"+licensorID {
			return &http.Response{
				StatusCode: 200,
				Body:       helptest.NoopCloser{bytes.NewBufferString(`{"invalid":"response"}`)},
				Header:     make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: 404,
			Body:       helptest.NoopCloser{bytes.NewBufferString("")},
			Header:     make(http.Header),
		}
	})

	client := NewClient(transport)

	_, err := client.Offer(vendorAPI, offerID)
	if err.Error() != "invalid offer" {
		t.Error("failed to return invalid offer error")
	}

	_, err = client.Licensor(vendorAPI, licensorID)
	if err.Error() != "invalid licensor" {
		t.Error("failed to return invalid licensor error")
	}
}

func TestInvalidJSON(t *testing.T) {
	vendorAPI := "https://licensezero.com"
	offerID := "186d34a9-c8f7-414c-91bc-a34b4553b91d"
	licensorID := "2d7aa023-8eca-4d15-ad43-6b15768a5293"
	transport := helptest.RoundTripFunc(func(req *http.Request) *http.Response {
		url := req.URL.String()
		if url == vendorAPI+"/offers/"+offerID ||
			url == vendorAPI+"/licensors/"+licensorID {
			return &http.Response{
				StatusCode: 200,
				Body:       helptest.NoopCloser{bytes.NewBufferString(`notvalidjson`)},
				Header:     make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: 404,
			Body:       helptest.NoopCloser{bytes.NewBufferString("")},
			Header:     make(http.Header),
		}
	})

	client := NewClient(transport)

	_, err := client.Offer(vendorAPI, offerID)
	if err.Error() != "invalid JSON" {
		t.Log(err.Error())
		t.Error("failed to return invalid offer error")
	}

	_, err = client.Licensor(vendorAPI, licensorID)
	if err.Error() != "invalid JSON" {
		t.Error("failed to return invalid licensor error")
	}
}
