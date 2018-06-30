package api

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "net/http"

type keyRequest struct {
	Action string `json:"action"`
}

type keyResponse struct {
	Error interface{} `json:"error"`
	Key   string      `json:"key"`
}

// FetchAgentPublicKey fetches the public Ed25519 signing key for licensezero.com.
func FetchAgentPublicKey() (string, error) {
	bodyData := keyRequest{Action: "key"}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return "", errors.New("error encoding agent key request body")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("error reading agent key response body")
	}
	var parsed keyResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return "", errors.New("error parsing agent key response body")
	}
	if message, ok := parsed.Error.(string); ok {
		return "", errors.New(message)
	}
	return parsed.Key, nil
}
