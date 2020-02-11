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
func ReadReceipts() ([]*api.Receipt, []error, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, nil, err
	}
	directoryPath := path.Join(configPath, "receipts")
	entries, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*api.Receipt{}, nil, nil
		}
		return nil, nil, err
	}
	var receipts []*api.Receipt
	var receiptErrors []error
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(configPath, "receipts", name)
		receipt, err := ReadReceipt(filePath)
		if err != nil {
			receiptErrors = append(receiptErrors, err)
			continue
		}
		receipts = append(receipts, receipt)
	}
	return receipts, receiptErrors, nil
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
func SaveReceipt(receipt *api.Receipt) error {
	json, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}
	err = os.MkdirAll(receiptsPath(configPath), 0755)
	if err != nil {
		return err
	}
	filePath, err := receiptPath(receipt, configPath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, json, 0644)
}

func receiptBasename(receipt *api.Receipt) (string, error) {
	digest := sha256.New()
	data, err := json.Marshal(receipt)
	if err != nil {
		return "", err
	}
	digest.Write(data)
	return hex.EncodeToString(digest.Sum(nil)), nil
}

// receiptPath calculates the file path for a receipt.
func receiptPath(receipt *api.Receipt, configPath string) (string, error) {
	basename, err := receiptBasename(receipt)
	if err != nil {
		return "", err
	}
	return path.Join(receiptsPath(configPath), basename+".json"), nil
}

func receiptsPath(configPath string) string {
	return path.Join(configPath, "receipts")
}
