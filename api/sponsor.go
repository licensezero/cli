package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "net/http"

type sponsorRequest struct {
	Action       string `json:"action"`
	ProjectID    string `json:"projectID"`
	Sponsor      string `json:"sponsor"`
	Jurisdiction string `json:"jurisdiction"`
	EMail        string `json:"email"`
}

type sponsorResponse struct {
	Error    interface{} `json:"error"`
	Location string      `json:"location"`
}

// Sponsor sends sponsor API requests.
func Sponsor(identity *data.Identity, projectID string) (string, error) {
	bodyData := sponsorRequest{
		Action:       "sponsor",
		ProjectID:    projectID,
		Sponsor:      identity.Name,
		Jurisdiction: identity.Jurisdiction,
		EMail:        identity.EMail,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return "", errors.New("could not construct sponsor request")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	if err != nil {
		return "", errors.New("error sending request")
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("invalid server response")
	}
	var parsed sponsorResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return "", err
	}
	if message, ok := parsed.Error.(string); ok {
		return "", errors.New(message)
	}
	location := parsed.Location
	return "https://licensezero.com" + location, nil
}
