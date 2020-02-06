package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Client responds to vendor API requests.
type Client interface {
	Offer(api string, offerID string) (*Offer, error)
	Licensor(api string, licensorID string) (*Licensor, error)
}

type httpClient struct {
	Client *http.Client
}

// NewClient reutrns a client using the given Transport.
func NewClient(t http.RoundTripper) Client {
	return &httpClient{Client: &http.Client{Transport: t}}
}

func (c *httpClient) Offer(api string, offerID string) (offer *Offer, err error) {
	response, err := c.Client.Get(api + "/offers/" + offerID)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var unstructured interface{}
	err = json.Unmarshal(body, &unstructured)
	if err != nil {
		return
	}
	offer, err = parseOffer(unstructured)
	if err != nil {
		return
	}
	return
}

func (c *httpClient) Licensor(api string, licensorID string) (licensor *Licensor, err error) {
	response, err := c.Client.Get(api + "/licensors/" + licensorID)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var unstructured interface{}
	err = json.Unmarshal(body, &unstructured)
	if err != nil {
		return
	}
	licensor, err = parseLicensor(unstructured)
	if err != nil {
		return
	}
	return
}
