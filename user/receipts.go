package user

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"os"
	"path"
)

// ReadReceipts reads all receipts in the configuration directory.
func ReadReceipts(configPath string) (receipts []*api.Receipt, errors []error, err error) {
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
func ReadReceipt(filePath string) (*api.Receipt, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var receipt api.Receipt
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

// SaveReceipt writes a receipt to the CLI configuration directory.
func SaveReceipt(receipt *api.Receipt, configPath string) error {
	json, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	err = os.MkdirAll(receiptsPath(configPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ReceiptPath(receipt, configPath), json, 0644)
}

func receiptBasename(api string, offerID string) string {
	digest := sha256.New()
	digest.Write([]byte(api + "/offers/" + offerID))
	return hex.EncodeToString(digest.Sum(nil))
}

// ReceiptPath calculates the file path for a receipt.
func ReceiptPath(receipt *api.Receipt, configPath string) string {
	basename := receiptBasename(
		receipt.License.Values.API,
		receipt.License.Values.OfferID,
	)
	return path.Join(receiptsPath(configPath), basename+".json")
}

func receiptsPath(configPath string) string {
	return path.Join(configPath, "receipts")
}
