package api

import "bytes"
import "encoding/json"
import "errors"
import "github.com/licensezero/cli/data"
import "net/http"

const TermsReference = "the terms of service at https://licensezero.com/terms/service"
const termsOfServiceStatement = "I agree to " + TermsReference + "."

type RegisterRequest struct {
	Action       string `json:"action"`
	Name         string `json:"licensee"`
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
		return errors.New("Server responded " + string(response.StatusCode))
	}
	return nil
}
