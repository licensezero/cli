package api

import "bytes"
import "encoding/json"
import "errors"
import "licensezero.com/cli/data"
import "io/ioutil"
import "net/http"
import "strconv"

type resetRequest struct {
	Action     string `json:"action"`
	LicensorID string `json:"licensorID"`
	EMail      string `json:"email"`
}

type resetResponse struct {
	Error interface{} `json:"error"`
}

// Reset sends reset API requests.
func Reset(identity *data.Identity, licensor *data.Licensor) error {
	bodyData := resetRequest{
		Action:     "reset",
		LicensorID: licensor.LicensorID,
		EMail:      identity.EMail,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return errors.New("could not construct reset request")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return errors.New("error sending request")
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("Server responded " + strconv.Itoa(response.StatusCode))
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var parsed resetResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return err
	}
	if message, ok := parsed.Error.(string); ok {
		return errors.New(message)
	}
	return nil
}
