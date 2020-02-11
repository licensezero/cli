package user

import (
	"errors"
	"fmt"
	"github.com/licensezero/helptest"
	"io/ioutil"
	"licensezero.com/licensezero/api"
	"os"
	"path"
	"testing"
)

func TestReadIdentity(t *testing.T) {
	directory, cleanup := helptest.TempDir(t, "licensezero")
	defer cleanup()
	email := "test@example.com"
	jurisdiction := "US-CA"
	name := "D Tester"
	err := ioutil.WriteFile(
		path.Join(directory, "identity.json"),
		[]byte("{\"email\": \""+email+"\", \"jurisdiction\": \""+jurisdiction+"\", \"name\": \""+name+"\"}"),
		0700,
	)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("LICENSEZERO_CONFIG", directory)

	result, err := ReadIdentity()
	if err != nil {
		t.Fatal("read error")
	}

	if result.Name != name {
		t.Error("did not read name")
	}
	if result.Jurisdiction != jurisdiction {
		t.Error("did not read jurisdiction")
	}
	if result.EMail != email {
		t.Error("did not read e-mail")
	}
}

func ExampleIdentity_ValidateReceipt_valid() {
	email := "test@example.com"
	jurisdiction := "US-CA"
	name := "Test Buyer"
	receipt := api.Receipt{
		License: api.License{
			Values: api.Values{
				Buyer: &api.Buyer{
					EMail:        email,
					Jurisdiction: jurisdiction,
					Name:         name,
				},
			},
		},
	}
	identity := Identity{
		EMail:        email,
		Jurisdiction: jurisdiction,
		Name:         name,
	}
	errs := identity.ValidateReceipt(&receipt)
	fmt.Println(errs == nil)
	// Output: true
}

func ExampleIdentity_ValidateReceipt_name() {
	email := "test@example.com"
	jurisdiction := "US-CA"
	name := "Test Buyer"
	receipt := api.Receipt{
		License: api.License{
			Values: api.Values{
				Buyer: &api.Buyer{
					EMail:        email,
					Jurisdiction: jurisdiction,
					Name:         "wrong",
				},
			},
		},
	}
	identity := Identity{
		EMail:        email,
		Jurisdiction: jurisdiction,
		Name:         name,
	}
	errs := identity.ValidateReceipt(&receipt)
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err.Error())
			fmt.Println(errors.Is(err, ErrNameMismatch))
		}
	}
	// Output:
	// name mismatch
	// true
}

func ExampleIdentity_ValidateReceipt_email() {
	email := "test@example.com"
	jurisdiction := "US-CA"
	name := "Test Buyer"
	receipt := api.Receipt{
		License: api.License{
			Values: api.Values{
				Buyer: &api.Buyer{
					EMail:        "wrong@example.com",
					Jurisdiction: jurisdiction,
					Name:         name,
				},
			},
		},
	}
	identity := Identity{
		EMail:        email,
		Jurisdiction: jurisdiction,
		Name:         name,
	}
	errs := identity.ValidateReceipt(&receipt)
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err.Error())
			fmt.Println(errors.Is(err, ErrEMailMismatch))
		}
	}
	// Output:
	// e-mail mismatch
	// true
}

func ExampleIdentity_ValidateReceipt_jurisdiction() {
	email := "test@example.com"
	jurisdiction := "US-CA"
	name := "Test Buyer"
	receipt := api.Receipt{
		License: api.License{
			Values: api.Values{
				Buyer: &api.Buyer{
					EMail:        email,
					Jurisdiction: "US-NV",
					Name:         name,
				},
			},
		},
	}
	identity := Identity{
		EMail:        email,
		Jurisdiction: jurisdiction,
		Name:         name,
	}
	errs := identity.ValidateReceipt(&receipt)
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err.Error())
			fmt.Println(errors.Is(err, ErrJurisdictionMismatch))
		}
	}
	// Output:
	// jurisdiction mismatch
	// true
}
