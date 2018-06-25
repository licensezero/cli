package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "net/http"

type PublicRequest struct {
	Action     string `json:"action"`
	LicensorID string `json:"licensorID"`
	Token      string `json:"token"`
	ProjectID  string `json:"projectID"`
	Terms      string `json:"terms"`
}

type PublicResponse struct {
	Error    interface{} `json:"error"`
	Metadata struct {
		License     string      `json:"license"`
		LicenseZero interface{} `json:"licensezero"`
	} `json:"metadata"`
	License struct {
		Document          string `json:"document"`
		LicensorSignature string `json:"licensorSignature"`
		AgentSignature    string `json:"agentSignature"`
	} `json:"license"`
}

func Public(licensor *data.Licensor, projectID string, terms string) (PublicResponse, error) {
	bodyData := PublicRequest{
		Action:     "public",
		ProjectID:  projectID,
		Terms:      terms,
		LicensorID: licensor.LicensorID,
		Token:      licensor.Token,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return PublicResponse{}, err
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return PublicResponse{}, err
	}
	var parsed PublicResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return PublicResponse{}, err
	}
	if message, ok := parsed.Error.(string); ok {
		return PublicResponse{}, errors.New(message)
	}
	return parsed, nil
}
