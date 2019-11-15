package api

import "bytes"
import "encoding/json"
import "errors"
import "licensezero.com/cli/data"
import "io/ioutil"
import "net/http"
import "strconv"

type repriceRequest struct {
	Action     string  `json:"action"`
	LicensorID string  `json:"licensorID"`
	Token      string  `json:"token"`
	ProjectID  string  `json:"offerID"`
	Pricing    Pricing `json:"pricing"`
}

type repriceResponse struct {
	Error interface{} `json:"error"`
}

// Reprice sends reprice API requests.
func Reprice(licensor *data.Licensor, offerID string, private, relicense uint) error {
	bodyData := repriceRequest{
		Action:     "reprice",
		LicensorID: licensor.LicensorID,
		ProjectID:  offerID,
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
	var parsed repriceResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return err
	}
	if message, ok := parsed.Error.(string); ok {
		return errors.New(message)
	}
	return nil
}
