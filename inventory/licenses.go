package inventory

type licenseType int

const (
	reciprocal licenseType = iota
	noncommercial
	unknown
)

func licenseTypeOf(identifier string) licenseType {
	switch identifier {
	case
		"Parity-1.0.0",
		"Parity-1.1.0",
		"Parity-2.0.0",
		"Parity-2.1.0",
		"Parity-2.2.0",
		"Parity-2.3.0",
		"Parity-2.3.1",
		"Parity-2.4.0",
		"Parity-3.0.0",
		"Parity-4.0.0",
		"Parity-5.0.0",
		"Parity-6.0.0",
		"Parity-7.0.0":
		return reciprocal
	case
		"Prosperity-1.0.0",
		"Prosperity-1.0.1",
		"Prosperity-1.1.0",
		"Prosperity-2.0.0",
		"Prosperity-3.0.0":
		return noncommercial
	default:
		return unknown
	}
}
