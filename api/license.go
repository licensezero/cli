package api

import "bytes"
import "encoding/json"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "net/http"

type LicenseRequest struct {
	Action     string `json:"action"`
	LicensorID string `json:"licensorID"`
	Token      string `json:"token"`
	ProjectID  string `json:"projectID"`
	Terms      string `json:"terms"`
}

type LicenseResponse struct {
	Metadata interface{} `json:"metadata"`
	License  struct {
		Document          string `json:"document"`
		LicensorSignature string `json:"licensorSignature"`
		AgentSignature    string `json:"agentSignature"`
	} `json:"license"`
}

func License(licensor *data.Licensor, projectID string, terms string) (LicenseResponse, error) {
	bodyData := LicenseRequest{
		Action:     "license",
		ProjectID:  projectID,
		Terms:      terms,
		LicensorID: licensor.LicensorID,
		Token:      licensor.Token,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return LicenseResponse{}, err
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return LicenseResponse{}, err
	}
	var parsed LicenseResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return LicenseResponse{}, err
	}
	return parsed, nil
}
