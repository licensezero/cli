package user

import (
	"bytes"
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
	Key       Key       `json:"key"`
	License   License   `json:"license"`
	Signature Signature `json:"signature"`
}

// Key represents a cryptographic signing public key.
type Key []byte

// Signature represents a cryptographic signature.
type Signature []byte

func toJSONHex(binary []byte) []byte {
	buffer := bytes.NewBufferString("\"")
	encoded := make([]byte, hex.EncodedLen(len(binary)))
	hex.Encode(encoded, binary)
	buffer.Write(encoded)
	buffer.WriteString("\"")
	return buffer.Bytes()
}

func fromJSONHex(data []byte) ([]byte, error) {
	var asString string
	err := json.Unmarshal(data, &asString)
	if err != nil {
		return nil, err
	}
	decoded := make([]byte, hex.DecodedLen(len(asString)))
	_, err = hex.Decode(decoded, []byte(asString))
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// MarshalJSON overrides default JSON serialization.
func (k *Key) MarshalJSON() ([]byte, error) {
	return toJSONHex(*k), nil
}

// UnmarshalJSON overrides default JSON serialization.
func (k *Key) UnmarshalJSON(data []byte) (err error) {
	bytes, err := fromJSONHex(data)
	if err != nil {
		return
	}
	*k = bytes
	return
}

// MarshalJSON overrides default JSON serialization.
func (s *Signature) MarshalJSON() ([]byte, error) {
	return toJSONHex(*s), nil
}

// UnmarshalJSON overrides default JSON serialization.
func (s *Signature) UnmarshalJSON(data []byte) (err error) {
	bytes, err := fromJSONHex(data)
	if err != nil {
		return
	}
	*s = bytes
	return
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
	API       string `json:"api"`
	Broker    Broker `json:"broker"`
	Buyer     Buyer  `json:"buyer"`
	Effective string `json:"effective"`
	Expires   string `json:"expires"`
	OfferID   string `json:"offerID"`
	OrderID   string `json:"orderID"`
	Price     Price  `json:"price"`
	Seller    Seller `json:"seller"`
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
	err = ValidateReceipt(&receipt)
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

// ValidateReceipt verifies that parsed JSON data conform
// to the JSON schema for receipt records.
func ValidateReceipt(receipt *Receipt) error {
	if receiptValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Receipt),
		)
		if err != nil {
			panic(err)
		}
		receiptValidator = schema
	}
	dataLoader := gojsonschema.NewGoLoader(receipt)
	result, err := receiptValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return ErrInvalidReceipt
	}
	return nil
}

// ValidateSignature validates the broker signature on a receipt.
func ValidateSignature(r *Receipt) error {
	serialized, err := json.Marshal(r.License)
	if err != nil {
		return err
	}
	return checkSignature(r.Key, r.Signature, serialized)
}

func checkSignature(publicKey []byte, signature []byte, json []byte) error {
	signatureValid := ed25519.Verify(
		publicKey,
		json,
		signature,
	)
	if !signatureValid {
		return errors.New("invalid signature")
	}
	return nil
}

// WriteReceipt writes a receipt to the CLI configuration directory.
func WriteReceipt(configPath string, receipt *Receipt) error {
	json, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	err = os.MkdirAll(receiptsPath(configPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(receiptPath(configPath, receipt), json, 0644)
}

func receiptBasename(api string, offerID string) string {
	digest := sha256.New()
	digest.Write([]byte(api + "/offers/" + offerID))
	return hex.EncodeToString(digest.Sum(nil))
}

// ReceiptPath calculates the file path for a receipt.
func receiptPath(configPath string, receipt *Receipt) string {
	basename := receiptBasename(
		receipt.License.Values.API,
		receipt.License.Values.OfferID,
	)
	return path.Join(receiptsPath(configPath), basename+".json")
}

func receiptsPath(configPath string) string {
	return path.Join(configPath, "receipts")
}
