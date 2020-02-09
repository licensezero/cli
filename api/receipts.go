package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/crypto/ed25519"
	"licensezero.com/licensezero/schemas"
)

// Receipt summarizes a receipt for a license.
type Receipt struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	KeyHex       string  `json:"key"`
	License      License `json:"license"`
	SignatureHex string  `json:"signature"`
}

// License holds data about the license a receipt documents.
type License struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	Form   string `json:"form"`
	Values Values `json:"values"`
}

// Values represent the values for blanks in the license form of a receipt.
type Values struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	API       string  `json:"api"`
	Broker    *Broker `json:"broker,omitempty"`
	Buyer     *Buyer  `json:"buyer"`
	Effective string  `json:"effective"`
	Expires   string  `json:"expires,omitempty"`
	OfferID   string  `json:"offerID"`
	OrderID   string  `json:"orderID"`
	Price     *Price  `json:"price,omitempty"`
	Seller    *Seller `json:"seller"`
	SellerID  string  `json:"sellerID"`
}

// Broker represents a party selling a license on behalf of a Seller.
type Broker struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
	Website      string `json:"website"`
}

// Buyer represents the party buying (or receiving) the license.
type Buyer struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
}

// ErrInvalidReceipt indicates that a Receipt does not conform
// to the JSON schema for receipt records.
var ErrInvalidReceipt = errors.New("invalid receipt")

var receiptValidator *gojsonschema.Schema

// Validate verifies that the receipt conforms
// to the JSON schema for receipt records.
func (receipt *Receipt) Validate() error {
	if receiptValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Receipt),
		)
		if err != nil {
			panic(err)
		}
		receiptValidator = schema
	}
	marshaled, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	dataLoader := gojsonschema.NewBytesLoader(marshaled)
	result, err := receiptValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return ErrInvalidReceipt
	}
	return nil
}

// ErrInvalidSignaure indicates that the siganture to a Receipt
// cannot be verified.
var ErrInvalidSignaure = errors.New("invalid signature")

// VerifySignature validates the broker signature on a receipt.
func (receipt *Receipt) VerifySignature() error {
	serialized, err := json.Marshal(receipt.License)
	if err != nil {
		return err
	}
	return checkSignature(receipt.KeyHex, receipt.SignatureHex, serialized)
}

func checkSignature(publicKey string, signature string, json []byte) error {
	signatureBytes := make([]byte, hex.DecodedLen(len(signature)))
	_, err := hex.Decode(signatureBytes, []byte(signature))
	if err != nil {
		return errors.New("bad signature")
	}
	publicKeyBytes := make([]byte, hex.DecodedLen(len(publicKey)))
	_, err = hex.Decode(publicKeyBytes, []byte(publicKey))
	if err != nil {
		return errors.New("bad public key")
	}
	signatureValid := ed25519.Verify(
		publicKeyBytes,
		json,
		signatureBytes,
	)
	if !signatureValid {
		return ErrInvalidSignaure
	}
	return nil
}
