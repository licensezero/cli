package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "io/ioutil"
import "net/http"
import "strconv"

const TermsReference = "the terms of service at https://licensezero.com/terms/service"
const termsOfServiceStatement = "I agree to " + TermsReference + "."

type RegisterRequest struct {
	Action       string `json:"action"`
	Name         string `json:"name"`
	Jurisdiction string `json:"jurisdiction"`
	EMail        string `json:"email"`
	Terms        string `json:"terms"`
}

func Register(identity *data.Identity) error {
	bodyData := RegisterRequest{
		Action:       "register",
		Name:         identity.Name,
		Jurisdiction: identity.Jurisdiction,
		EMail:        identity.EMail,
		Terms:        termsOfServiceStatement,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return errors.New("could not construct register request")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("Server responded " + strconv.Itoa(response.StatusCode))
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var parsed struct {
		Error interface{} `json:"error"`
	}
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return err
	}
	if message, ok := parsed.Error.(string); ok {
		return errors.New(message)
	}
	return nil
}
