package api

import "bytes"
import "encoding/json"
import "errors"
import "licensezero.com/cli/data"
import "io/ioutil"
import "net/http"
import "strconv"

// TermsReference includes the text required in terms of service agreement statements to the API.
const TermsReference = "the terms of service at https://licensezero.com/terms/service"
const termsOfServiceStatement = "I agree to " + TermsReference + "."

type registerRequest struct {
	Action       string `json:"action"`
	Name         string `json:"name"`
	Jurisdiction string `json:"jurisdiction"`
	EMail        string `json:"email"`
	Terms        string `json:"terms"`
}

type registerResponse struct {
	Error interface{} `json:"error"`
}

// Register sends a register API request.
func Register(identity *data.Identity) error {
	bodyData := registerRequest{
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
	var parsed registerResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return err
	}
	if message, ok := parsed.Error.(string); ok {
		return errors.New(message)
	}
	return nil
}
