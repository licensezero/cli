package user

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"licensezero.com/licensezero/schemas"
	"os"
	"path"
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

// Price represents the price paid for the license.
type Price struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	Amount   uint   `json:"amount"`
	Currency string `json:"currency"`
}

// Seller represents the party selling the license.
// Seller is usually the developer of the software being licensed.
type Seller struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
	SellerID     string `json:"sellerID"`
}

// ReadReceipts reads all receipts in the configuration directory.
func ReadReceipts(configPath string) (receipts []*Receipt, errors []error, err error) {
	directoryPath := path.Join(configPath, "receipts")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return
		}
		return nil, nil, directoryReadError
	}
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(configPath, "receipts", name)
		receipt, err := ReadReceipt(filePath)
		if err != nil {
			errors = append(errors, err)
		} else {
			receipts = append(receipts, receipt)
		}
	}
	return
}

// ReadReceipt reads a receipt record from a file.
func ReadReceipt(filePath string) (*Receipt, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var receipt Receipt
	err = json.Unmarshal(data, &receipt)
	if err != nil {
		return nil, err
	}
	err = receipt.Validate()
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

var receiptValidator *gojsonschema.Schema = nil

// ErrInvalidReceipt indicates that parsing or validation
// failed because the data do not conform to the receipt
// JSON schema.
var ErrInvalidReceipt = errors.New("invalid receipt")

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
		return errors.New("invalid signature")
	}
	publicKeyBytes := make([]byte, hex.DecodedLen(len(publicKey)))
	_, err = hex.Decode(publicKeyBytes, []byte(publicKey))
	if err != nil {
		return errors.New("invalid public key")
	}
	signatureValid := ed25519.Verify(
		publicKeyBytes,
		json,
		signatureBytes,
	)
	if !signatureValid {
		return errors.New("invalid signature")
	}
	return nil
}

// Save writes a receipt to the CLI configuration directory.
func (receipt *Receipt) Save(configPath string) error {
	json, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	err = os.MkdirAll(receiptsPath(configPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(receipt.Path(configPath), json, 0644)
}

func receiptBasename(api string, offerID string) string {
	digest := sha256.New()
	digest.Write([]byte(api + "/offers/" + offerID))
	return hex.EncodeToString(digest.Sum(nil))
}

// Path calculates the file path for a receipt.
func (receipt *Receipt) Path(configPath string) string {
	basename := receiptBasename(
		receipt.License.Values.API,
		receipt.License.Values.OfferID,
	)
	return path.Join(receiptsPath(configPath), basename+".json")
}

func receiptsPath(configPath string) string {
	return path.Join(configPath, "receipts")
}
