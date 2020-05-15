package api

import "bytes"
import "encoding/json"
import "errors"
import "licensezero.com/cli/data"
import "io/ioutil"
import "net/http"

type publicRequest struct {
	Action      string `json:"action"`
	DeveloperID string `json:"developerID"`
	Token       string `json:"token"`
	OfferID     string `json:"offerID"`
	Terms       string `json:"terms"`
}

// PublicResponse contains API instructions for adding public licensing information.
type PublicResponse struct {
	Error    interface{} `json:"error"`
	Metadata struct {
		License     string      `json:"license"`
		LicenseZero interface{} `json:"licensezero"`
	} `json:"metadata"`
	License struct {
		Document           string `json:"document"`
		DeveloperSignature string `json:"developerSignature"`
		AgentSignature     string `json:"agentSignature"`
	} `json:"license"`
}

// Public sends a public API request.
func Public(developer *data.Developer, offerID string, terms string) (*PublicResponse, error) {
	bodyData := publicRequest{
		Action:      "public",
		OfferID:     offerID,
		Terms:       terms,
		DeveloperID: developer.DeveloperID,
		Token:       developer.Token,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("error sending request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var parsed PublicResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, err
	}
	if message, ok := parsed.Error.(string); ok {
		return nil, errors.New(message)
	}
	return &parsed, nil
}
