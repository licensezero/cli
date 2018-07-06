package api

import "encoding/json"
import "strings"
import "testing"

func TestRepriceMissing(t *testing.T) {
	data := Pricing{
		Private: 100,
	}
	json, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(string(json), "\"relicense\"") {
		t.Error("contains relicense property")
	}
}
