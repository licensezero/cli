package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "net/http"

type RepriceRequest struct {
	Action     string  `json:"action"`
	LicensorID string  `json:"licensorID"`
	Token      string  `json:"token"`
	ProjectID  string  `json:"projectID"`
	Pricing    Pricing `json:"pricing"`
}

func Reprice(licensor *data.Licensor, projectID string, private, relicense uint) error {
	bodyData := RepriceRequest{
		Action:     "reprice",
		LicensorID: licensor.LicensorID,
		Token:      licensor.Token,
		Pricing: Pricing{
			Private:   private,
			Relicense: relicense,
		},
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("Server responded " + string(response.StatusCode))
	}
	return nil
}
