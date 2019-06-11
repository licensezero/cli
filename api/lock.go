package api

import "bytes"
import "encoding/json"
import "errors"
import "licensezero.com/cli/data"
import "io/ioutil"
import "net/http"
import "strconv"

import "fmt"

type lockRequest struct {
	Action     string `json:"action"`
	LicensorID string `json:"licensorID"`
	Token      string `json:"token"`
	ProjectID  string `json:"projectID"`
	Unlock     string `json:"unlock"`
}

type lockResponse struct {
	Error interface{} `json:"error"`
}

// Lock sends a lock API request.
func Lock(licensor *data.Licensor, projectID string, unlock string) error {
	bodyData := lockRequest{
		Action:     "lock",
		ProjectID:  projectID,
		Unlock:     unlock,
		LicensorID: licensor.LicensorID,
		Token:      licensor.Token,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("Server responded " + strconv.Itoa(response.StatusCode))
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var parsed lockResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return err
	}
	if message, ok := parsed.Error.(string); ok {
		return errors.New(message)
	}
	return nil
}
