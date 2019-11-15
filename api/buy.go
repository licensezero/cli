package api

import "bytes"
import "encoding/json"
import "errors"
import "licensezero.com/cli/data"
import "io/ioutil"
import "net/http"

type buyRequest struct {
	Action       string   `json:"action"`
	Projects     []string `json:"projects"`
	Name         string   `json:"licensee"`
	Jurisdiction string   `json:"jurisdiction"`
	EMail        string   `json:"email"`
	Person       string   `json:"person"`
}

type buyResponse struct {
	Error    interface{} `json:"error"`
	Location string      `json:"location"`
}

// Buy sends a buy API request.
func Buy(identity *data.Identity, offerIDs []string) (string, error) {
	bodyData := buyRequest{
		Action:       "order",
		Projects:     offerIDs,
		Name:         identity.Name,
		Jurisdiction: identity.Jurisdiction,
		EMail:        identity.EMail,
		Person:       "I am a person, not a legal entity.",
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return "", errors.New("could not construct quote request")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("error sending request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("invalid server response")
	}
	var parsed buyResponse
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
