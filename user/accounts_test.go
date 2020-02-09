package user

import (
	"github.com/licensezero/helptest"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestReadAccounts(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	receipts := path.Join(directory, "accounts")
	err := os.MkdirAll(receipts, 0700)
	if err != nil {
		t.Fatal(err)
	}

	account := path.Join(receipts, "first.json")
	api := "https://broker.licensezero.com"
	sellerID := "71ea37d7-6a1a-4072-a64b-84d0236edfe6"
	token := "xxxxxx"
	err = ioutil.WriteFile(
		account,
		[]byte("{\"api\": \""+api+"\", \"sellerID\": \""+sellerID+"\", \"token\": \""+token+"\"}"),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	results, err := ReadAccounts(directory)
	if err != nil {
		t.Fatal("read error")
	}

	if len(results) != 1 {
		t.Fatal("did not find one account")
	}
	result := results[0]
	if result.API != api {
		t.Error("did not read API")
	}
	if result.SellerID != sellerID {
		t.Error("did not read seller ID")
	}
}
