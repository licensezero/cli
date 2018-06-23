package api

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "net/http"

type ProjectRequest struct {
	Action    string `json:"action"`
	ProjectID string `json:"projectID"`
}

type ProjectResponse struct {
	Licensor LicensorInformation `json:"licensor"`
}

func Project(projectID string) (ProjectResponse, error) {
	bodyData := ProjectRequest{
		Action:    "project",
		ProjectID: projectID,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return ProjectResponse{}, errors.New("error encoding agent key request body")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ProjectResponse{}, errors.New("error reading agent key response body")
	}
	var parsed ProjectResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return ProjectResponse{}, errors.New("error parsing agent key response body")
	}
	return parsed, nil
}
