package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// GetKey fetches the public signing key of a vendor API.
func GetKey(api string) (string, error) {
	response, err := http.Get(api + "/key")
	if err != nil {
		return "", errors.New("error sending request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("error reading response body")
	}
	var parsed string
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return "", errors.New("error parsing response body")
	}
	return parsed, nil
}
