package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "net/http"
import "strconv"

type LockRequest struct {
	Action     string `json:"action"`
	LicensorID string `json:"licensorID"`
	Token      string `json:"token"`
	ProjectID  string `json:"projectID"`
	Unlock     string `json:"unloack"`
}

func Lock(licensor *data.Licensor, projectID string, unlock string) error {
	bodyData := LockRequest{
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
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("Server responded " + strconv.Itoa(response.StatusCode))
	}
	return nil
}
