package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "net/http"
import "strconv"

type WaiveRequest struct {
	Action       string `json:"action"`
	LicensorID   string `json:"licensorID"`
	Token        string `json:"token"`
	ProjectID    string `json:"projectID"`
	Beneficiary  string `json:"beneficiary"`
	Jurisdiction string `json:"jurisdiction"`
	Term         string `json:"term"`
}

func Waive(licensor *data.Licensor, projectID, beneficiary, jurisdiction, term string) ([]byte, error) {
	bodyData := WaiveRequest{
		Action:       "waive",
		LicensorID:   licensor.LicensorID,
		Token:        licensor.Token,
		ProjectID:    projectID,
		Beneficiary:  beneficiary,
		Jurisdiction: jurisdiction,
		Term:         term,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, errors.New("Error serializing request body.")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("Server responded " + strconv.Itoa(response.StatusCode))
	}
	return ioutil.ReadAll(response.Body)
}
