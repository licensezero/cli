package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Licensor represents data about a licensor selling through a vendor.
type Licensor struct {
	Name         string
	Jurisdiction string
}

// GetLicensor fetches information abourt a licensor from a vendor API.
func GetLicensor(api string, licensorID string) (*Licensor, error) {
	response, err := http.Get(api + "/licensors/" + licensorID)
	if err != nil {
		return nil, errors.New("error sending request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}
	var parsed Licensor
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, errors.New("error parsing response body")
	}
	return &parsed, nil
}
