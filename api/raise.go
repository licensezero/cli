package api

import "bytes"
import "encoding/json"
import "errors"
import "licensezero.com/cli/data"
import "io/ioutil"
import "fmt"
import "net/http"
import "strconv"

type raiseRequest struct {
	Action      string `json:"action"`
	DeveloperID string `json:"developerID"`
	Token       string `json:"token"`
	OfferID     string `json:"offerID"`
	Commission  uint   `json:"commission"`
}

type raiseResponse struct {
	Error interface{} `json:"error"`
}

// Raise sends raise API requests.
func Raise(developer *data.Developer, offerID string, commission uint) error {
	bodyData := raiseRequest{
		Action:      "raise",
		DeveloperID: developer.DeveloperID,
		OfferID:     offerID,
		Token:       developer.Token,
		Commission:  commission,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}
	fmt.Println(bodyData)
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
