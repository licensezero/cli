package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "net/http"
import "strconv"

type RetractRequest struct {
	Action     string `json:"action"`
	LicensorID string `json:"licensorID"`
	Token      string `json:"token"`
	ProjectID  string `json:"projectID"`
}

func Retract(licensor *data.Licensor, projectID string) error {
	bodyData := RetractRequest{
		Action:     "retract",
		LicensorID: licensor.LicensorID,
		Token:      licensor.Token,
		ProjectID:  projectID,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("Server responded " + strconv.Itoa(response.StatusCode))
	}
	return nil
}
