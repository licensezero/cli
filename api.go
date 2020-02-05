package cli

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// TODO: Figure out how to mock API responses.

func getOffer(api string, offerID string) (offer *Offer, err error) {
	response, err := http.Get(api + "/offers/" + offerID)
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
