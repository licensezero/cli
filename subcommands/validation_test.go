package subcommands

import (
	"testing"
)

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

func TestValidUUIDv4(t *testing.T) {
	if !validID("1f838e1d-c98f-44a3-a4e8-15267a0f0777") {
		t.Error("valid UUIDv4 failed")
	}
	if !validID("0424944d-a682-4301-8d7d-3a9a4173be48") {
		t.Error("valid UUIDv4 failed")
	}
	if validID("0424944d-a682-5301-8d7d-3a9a4173be48") {
		//                      ^
		t.Error("invalid UUIDv4 passed")
	}
	if validID("0424944d-a682-4301-cd7d-3a9a4173be48") {
		//                           ^
		t.Error("invalid UUIDv4 passed")
	}
	if validID("0424944d-a682-4301-8d7d-3a9a4173be48abab") {
		//                                            ^
		t.Error("invalid UUIDv4 passed")
	}
}
