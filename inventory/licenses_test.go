package inventory

import (
	"testing"
)

func TestLicenseTypeOf(t *testing.T) {
	cases := map[string]licenseType{
		"Parity-7.0.0":     reciprocal,
		"Prosperity-3.0.0": noncommercial,
		"MIT":              unknown,
	}
	for public, licenseType := range cases {
		result := licenseTypeOf(public)
		if result != licenseType {
			t.Error(public)
		}
	}
}
