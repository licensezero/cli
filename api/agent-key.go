package api

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "net/http"

type KeyRequest struct {
	Action string `json:"action"`
}

type KeyResponse struct {
	Key string `json:"key"`
}

func FetchAgentPublicKey() (string, error) {
	bodyData := KeyRequest{Action: "key"}
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
	var parsed KeyResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return "", errors.New("error parsing agent key response body")
	}
	return parsed.Key, nil
}
