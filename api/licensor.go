package api

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "net/http"

type licensorRequest struct {
	Action     string `json:"action"`
	LicensorID string `json:"licensorID"`
}

// ProjectInformation describes information on a project from an API Licensor request.
type ProjectInformation struct {
	OfferID   string `json:"offerID"`
	Offered   string `json:"offered"`
	Retracted string `json:"retracted,omitempty"`
}

type licensorResponse struct {
	Error        interface{}          `json:"error"`
	Name         string               `json:"name"`
	Jurisdiction string               `json:"jurisdiction"`
	PublicKey    string               `json:"publicKey"`
	Projects     []ProjectInformation `json:"projects"`
}

// Licensor sends a licensor API request.
func Licensor(licensorID string) (*LicensorInformation, []ProjectInformation, error) {
	bodyData := licensorRequest{
		Action:     "licensor",
		LicensorID: licensorID,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, nil, errors.New("error encoding licensor request body")
	}
	response, err := http.Post("https://licensezero.com/api/v0", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, errors.New("error sending licensor request")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, errors.New("error reading licensor response body")
	}
	var parsed licensorResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, nil, errors.New("error parsing licensor response body")
	}
	if message, ok := parsed.Error.(string); ok {
		return nil, nil, errors.New(message)
	}
	licensor := LicensorInformation{
		Name:         parsed.Name,
		Jurisdiction: parsed.Jurisdiction,
		PublicKey:    parsed.PublicKey,
	}
	return &licensor, parsed.Projects, nil
}
