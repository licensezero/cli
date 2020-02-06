package user

import (
	"bytes"
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
	API       string
	Effective string
	Expires   string
	Key       string
	OfferID   string
	OrderID   string
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
		receipt, err := readReceipt(filePath)
		if err != nil {
			errors = append(errors, err)
		} else {
			receipts = append(receipts, receipt)
		}
	}
	return
}

func readReceipt(filePath string) (*Receipt, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var unstructured interface{}
	err = json.Unmarshal(data, &unstructured)
	if err != nil {
		return nil, err
	}
	err = validateReceipt(unstructured)
	if err != nil {
		return nil, err
	}
	return parseReceipt(unstructured), nil
}

var receiptValidator *gojsonschema.Schema = nil

func validateReceipt(unstructured interface{}) error {
	if receiptValidator == nil {
		schema, err := schemas.Loader().Compile(
			gojsonschema.NewStringLoader(schemas.Receipt),
		)
		if err != nil {
			panic(err)
		}
		receiptValidator = schema
	}
	dataLoader := gojsonschema.NewGoLoader(unstructured)
	result, err := receiptValidator.Validate(dataLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return errors.New("invalid receipt")
	}
	return nil
}

func parseReceipt(unstructured interface{}) (r *Receipt) {
	asMap := unstructured.(map[string]interface{})
	license := asMap["license"].(map[string]interface{})
	values := license["values"].(map[string]interface{})
	expires, _ := values["expires"].(string)
	return &Receipt{
		API:       values["api"].(string),
		Key:       asMap["key"].(string),
		OfferID:   values["offerID"].(string),
		OrderID:   values["orderID"].(string),
		Effective: values["effective"].(string),
		Expires:   expires,
	}
}

func validateSignature(unstructured interface{}) error {
	serialized, err := serializeReceiptLicense(unstructured)
	if err != nil {
		return err
	}
	asMap := unstructured.(map[string]interface{})
	key := asMap["key"].(string)
	signature := asMap["signature"].(string)
	return checkSignature(key, signature, serialized)
}

func serializeReceiptLicense(unstructured interface{}) ([]byte, error) {
	asMap := unstructured.(map[string]interface{})
	license := asMap["license"].(map[string]interface{})
	serialized, err := json.Marshal(license)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer([]byte{})
	err = json.Compact(buffer, serialized)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
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
