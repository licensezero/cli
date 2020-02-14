package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

// BrokerServer responds to broker server requests.
type BrokerServer struct {
	Client *http.Client
	Base   string
}

// Offer requests information about a license offer.
func (b *BrokerServer) Offer(offerID string) (offer *Offer, err error) {
	response, err := b.Client.Get(b.Base + "/offers/" + offerID)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &offer)
	if err != nil {
		return nil, errors.New("invalid JSON")
	}
	err = offer.Validate()
	if err != nil {
		return nil, errors.New("invalid offer")
	}
	return
}

// Seller requests information about a license seller.
func (b *BrokerServer) Seller(sellerID string) (seller *Seller, err error) {
	response, err := b.Client.Get(b.Base + "/sellers/" + sellerID)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &seller)
	if err != nil {
		return nil, errors.New("invalid JSON")
	}
	err = seller.Validate()
	if err != nil {
		return nil, errors.New("invalid seller")
	}
	return
}

// Register gets the broker's signing key register.
func (b *BrokerServer) Register() (register *Register, err error) {
	response, err := b.Client.Get(b.Base + "/register")
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &register)
	if err != nil {
		return nil, errors.New("invalid JSON")
	}
	err = register.Validate()
	if err != nil {
		return nil, errors.New("invalid register")
	}
	return
}

// Latest gets the latest receipt for an order.
func (b *BrokerServer) Latest(orderID string) (receipt *Receipt, err error) {
	response, err := b.Client.Get(b.Base + "/orders/" + orderID + "/latest")
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &receipt)
	if err != nil {
		return nil, errors.New("invalid JSON")
	}
	return
}

// Broker gets information about the broker operating the server.
func (b *BrokerServer) Broker() (broker *Broker, err error) {
	response, err := b.Client.Get(b.Base + "/broker")
	if err != nil {
		return
	}
	if response.StatusCode == 404 {
		return nil, nil
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &broker)
	if err != nil {
		return nil, errors.New("invalid JSON")
	}
	err = broker.Validate()
	if err != nil {
		return nil, errors.New("invalid broker")
	}
	return
}

// Order creates an order for licenses.
func (b *BrokerServer) Order(
	email, jurisdiction, name string,
	offerIDs []string,
) (string, error) {
	var buffer bytes.Buffer
	postBody := multipart.NewWriter(&buffer)
	postBody.WriteField("email", email)
	postBody.WriteField("name", name)
	postBody.WriteField("jurisdiction", jurisdiction)
	for _, offerID := range offerIDs {
		postBody.WriteField("offerIDs[]", offerID)
	}
	postBody.Close()
	request, err := http.NewRequest("POST", b.Base+"/buy", &buffer)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", postBody.FormDataContentType())
	response, err := b.Client.Do(request)
	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("bad status: %s", response.Status)
	}
	location, err := response.Location()
	if err != nil {
		return "", err
	}
	return location.String(), nil
}

// RegisterSeller registers a seller to sell through a broker.
func (b *BrokerServer) RegisterSeller(
	email, jurisdiction, name string,
) error {
	var buffer bytes.Buffer
	postBody := multipart.NewWriter(&buffer)
	postBody.WriteField("email", email)
	postBody.WriteField("name", name)
	postBody.WriteField("jurisdiction", jurisdiction)
	postBody.Close()
	request, err := http.NewRequest("POST", b.Base+"/sellers", &buffer)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", postBody.FormDataContentType())
	response, err := b.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 201 {
		return fmt.Errorf("bad status: %s", response.Status)
	}
	return nil
}

// ResetToken requests a seller token reset.
func (b *BrokerServer) ResetToken(sellerID string) error {
	var buffer bytes.Buffer
	postBody := multipart.NewWriter(&buffer)
	postBody.WriteField("sellerID", sellerID)
	postBody.Close()
	request, err := http.NewRequest("POST", b.Base+"/reset", &buffer)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", postBody.FormDataContentType())
	response, err := b.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 201 {
		return fmt.Errorf("bad status: %s", response.Status)
	}
	return nil
}

// RegisterOffer offers licenses for sale.
func (b *BrokerServer) RegisterOffer(
	sellerID string,
	token string,
	repository string,
	description string,
	price uint,
	relicense *uint,
) (string, error) {
	var buffer bytes.Buffer
	postBody := multipart.NewWriter(&buffer)
	postBody.WriteField("sellerID", sellerID)
	postBody.WriteField("token", token)
	postBody.WriteField("repository", repository)
	postBody.WriteField("description", description)
	postBody.WriteField("price", formatPrice(price))
	if relicense != nil {
		postBody.WriteField("relicense", formatPrice(*relicense))
	}
	postBody.Close()
	request, err := http.NewRequest("POST", b.Base+"/offers", &buffer)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", postBody.FormDataContentType())
	response, err := b.Client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("bad status: %s", response.Status)
	}
	location, err := response.Location()
	if err != nil {
		return "", err
	}
	return location.String(), nil
}

func formatPrice(price uint) string {
	return strconv.FormatUint(uint64(price), 10)
}

// Reprice offers licenses for sale.
func (b *BrokerServer) Reprice(
	sellerID string,
	token string,
	offerID string,
	price uint,
	relicense *uint,
) error {
	var buffer bytes.Buffer
	postBody := multipart.NewWriter(&buffer)
	postBody.WriteField("sellerID", sellerID)
	postBody.WriteField("token", token)
	postBody.WriteField("offerID", offerID)
	postBody.WriteField("price", formatPrice(price))
	if relicense != nil {
		postBody.WriteField("relicense", formatPrice(*relicense))
	}
	postBody.Close()
	request, err := http.NewRequest("PATCH", b.Base+"/offers/"+offerID, &buffer)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", postBody.FormDataContentType())
	response, err := b.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("bad status: %s", response.Status)
	}
	return nil
}

// Lock locks license pricing and availability.
func (b *BrokerServer) Lock(
	sellerID string,
	token string,
	offerID string,
	unlock string,
) error {
	var buffer bytes.Buffer
	postBody := multipart.NewWriter(&buffer)
	postBody.WriteField("sellerID", sellerID)
	postBody.WriteField("token", token)
	postBody.WriteField("offerID", offerID)
	postBody.WriteField("unlock", unlock)
	postBody.Close()
	request, err := http.NewRequest("PATCH", b.Base+"/offers/"+offerID, &buffer)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", postBody.FormDataContentType())
	response, err := b.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("bad status: %s", response.Status)
	}
	return nil
}
