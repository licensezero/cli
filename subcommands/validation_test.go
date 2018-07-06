package subcommands

import "testing"

func TestValidJurisdiction(t *testing.T) {
	if !validJurisdiction("US-CA") {
		t.Error("US-CA fails")
	}
	if !validJurisdiction("US-TX") {
		t.Error("US-CA fails")
	}
	if !validJurisdiction("RU-MOS") {
		t.Error("RU-MOS fails")
	}
}
