package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// BrokerServer responds to broker API requests.
type BrokerServer struct {
	Client *http.Client
	Base   string
}

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

func (b *BrokerServer) Register() (register *Register, err error) {
	response, err := b.Client.Get(b.Base + "/keys")
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
		return nil, errors.New("invalid keys")
	}
	return
}
