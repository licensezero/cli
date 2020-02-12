package api

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func ExampleRegister_ValidateEffectiveDate_valid() {
	keyHex := strings.Repeat("a", 64)
	receipt := makeReceipt(keyHex, "2019-01-01T00:00:00Z")
	from, _ := time.Parse(time.RFC3339, "2018-01-01T00:00:00Z")
	keys := make(map[string]Timeframe)
	keys[keyHex] = Timeframe{From: &RegisterTime{from}}
	register := Register{
		Updated: "2020-01-01T00:00:00Z",
		Keys:    keys,
	}
	err := register.ValidateEffectiveDate(&receipt)
	fmt.Println(err == nil)
	// Output: true
}

func ExampleRegister_ValidateEffectiveDate_backdated() {
	keyHex := strings.Repeat("a", 64)
	receipt := makeReceipt(keyHex, "2018-01-01T00:00:00Z")
	from, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	keys := make(map[string]Timeframe)
	keys[keyHex] = Timeframe{From: &RegisterTime{from}}
	register := Register{
		Updated: "2020-01-01T00:00:00Z",
		Keys:    keys,
	}
	err := register.ValidateEffectiveDate(&receipt)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Output: backdated signature
}

func ExampleRegister_ValidateEffectiveDate_postdated() {
	keyHex := strings.Repeat("a", 64)
	from, _ := time.Parse(time.RFC3339, "2017-01-01T00:00:00Z")
	through, _ := time.Parse(time.RFC3339, "2018-01-01T00:00:00Z")
	receipt := makeReceipt(keyHex, "2019-01-01T00:00:00Z")
	keys := make(map[string]Timeframe)
	keys[keyHex] = Timeframe{
		From:    &RegisterTime{from},
		Through: &RegisterTime{through},
	}
	register := Register{
		Updated: "2020-01-01T00:00:00Z",
		Keys:    keys,
	}
	err := register.ValidateEffectiveDate(&receipt)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Output: postdated signature
}

func makeReceipt(keyHex string, effective string) Receipt {
	return Receipt{
		SignatureHex: "",
		KeyHex:       keyHex,
		License: License{
			Form: "test form",
			Values: Values{
				Server:    "https://broker.licensezero.com",
				Effective: effective,
				OfferID:   uuid.New().String(),
				OrderID:   uuid.New().String(),
				Buyer: &Buyer{
					Name:         "Buyer",
					EMail:        "buyer@example.com",
					Jurisdiction: "US-TX",
				},
				Seller: &Seller{
					Name:         "Seller",
					EMail:        "seller@example.com",
					Jurisdiction: "US-CA",
				},
			},
		},
	}
}
