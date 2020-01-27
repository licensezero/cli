package data

import "bytes"
import "encoding/json"
import "errors"
import "io/ioutil"
import "os"
import "path"

// Receipt describes the data in a receipt record.
type Receipt struct {
	Schema    string         `json:"schema"`
	Key       string         `json:"key"`
	Signature string         `json:"signature"`
	License   ReceiptLicense `json:"license"`
}

// ReceiptLicense describes a receipt's data about the license granted.
type ReceiptLicense struct {
	Values ReceiptLicenseValues `json:"values"`
	Form   string               `json:"form"`
}

// ReceiptLicenseValues describes the values plugged into a license form.
type ReceiptLicenseValues struct {
	Order     string          `json:"order"`
	Offer     string          `json:"offer"`
	Effective string          `json:"effective"`
	Price     Price           `json:"price"`
	Expires   string          `json:"expires"`
	Licensee  ReceiptLicensee `json:"licensee"`
	Licensor  ReceiptLicensor `json:"licensor"`
	Vendor    ReceiptVendor   `json:"vendor"`
}

// Price describes a particular price in a particular currency.
type Price struct {
	Amount   uint   `json:"amount"`
	Currency string `json:"currency"`
}

// ReceiptLicensee describes the one receiving a license.
type ReceiptLicensee struct {
	EMail        string `json:"email"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
}

// ReceiptLicensor describves the one giving a license.
type ReceiptLicensor struct {
	EMail        string `json:"email"`
	ID           string `json:"id"`
	Jurisdiction string `json:"jurisdiction"`
	Name         string `json:"name"`
}

// ReceiptVendor describes the one selling a license.
type ReceiptVendor struct {
	API      string `json:"api"`
	EMail    string `json:"email"`
	Homepage string `json:"homepage"`
	Name     string `json:"name"`
}

func receiptPath(home string, receipt *Receipt) string {
	return path.Join(
		receiptsPath(home),
		receipt.License.Values.Vendor.API,
		receipt.License.Values.Offer+".json",
	)
}

func receiptsPath(home string) string {
	return path.Join(ConfigPath(home), "receipts")
}

// ReadReceipts reads all saved receipts from the CLI configuration directory.
func ReadReceipts(home string) ([]Receipt, error) {
	directoryPath := path.Join(ConfigPath(home), "receipts")
	entries, directoryReadError := ioutil.ReadDir(directoryPath)
	if directoryReadError != nil {
		if os.IsNotExist(directoryReadError) {
			return []Receipt{}, nil
		}
		return nil, directoryReadError
	}
	var returned []Receipt
	for _, entry := range entries {
		name := entry.Name()
		filePath := path.Join(home, "receipts", name)
		receipt, err := ReadReceipt(filePath)
		if err != nil {
			return nil, err
		}
		returned = append(returned, *receipt)
	}
	return returned, nil
}

// ReadReceipt reads a receipt file from disk.
func ReadReceipt(filePath string) (*Receipt, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var file Receipt
	err = json.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// WriteReceipt writes a receipt file to the CLI configuration directory.
func WriteReceipt(home string, receipt *Receipt) error {
	json, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	err = os.MkdirAll(receiptsPath(home), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(receiptPath(home, receipt), json, 064)
}

// CheckReceiptSignature verifies the signatures to a liecnse envelope.
func CheckReceiptSignature(receipt *Receipt, publicKey string) error {
	serialized, err := json.Marshal(receipt.License)
	if err != nil {
		return errors.New("could not serialize receipt license")
	}
	compacted := bytes.NewBuffer([]byte{})
	err = json.Compact(compacted, serialized)
	if err != nil {
		return errors.New("could not compact receipt license")
	}
	err = checkSignature(
		publicKey,
		receipt.Signature,
		compacted.Bytes(),
	)
	if err != nil {
		return err
	}
	return nil
}
