package cli

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestReadAccounts(t *testing.T) {
	withTestDir(t, func(directory string) {
		receipts := path.Join(directory, "accounts")
		err := os.MkdirAll(receipts, 0700)
		if err != nil {
			t.Fatal(err)
		}

		account := path.Join(receipts, "first.json")
		api := "https://api.commonform.com"
		licensorID := "71ea37d7-6a1a-4072-a64b-84d0236edfe6"
		token := "xxxxxx"
		err = ioutil.WriteFile(
			account,
			[]byte("{\"api\": \""+api+"\", \"licensorID\": \""+licensorID+"\", \"token\": \""+token+"\"}"),
			0700,
		)
		if err != nil {
			t.Fatal(err)
		}

		results, err := readAccounts(directory)
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
		if result.LicensorID != licensorID {
			t.Error("did not read licensor ID")
		}
	})
}
