package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "net/http"

type BuyRequest struct {
	Action       string   `json:"action"`
	Projects     []string `json:"projects"`
	Name         string   `json:"licensee"`
	Jurisdiction string   `json:"jurisdiction"`
	EMail        string   `json:"email"`
	Person       string   `json:"person"`
}

type BuyResponse struct {
	Error    interface{} `json:"error"`
	Location string      `json:"location"`
}

func Buy(identity *data.Identity, projectIDs []string) (string, error) {
	bodyData := BuyRequest{
		Action:       "order",
		Projects:     projectIDs,
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
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("invalid server response")
	}
	var parsed BuyResponse
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
