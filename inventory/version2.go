package inventory

// Version2Envelope describes a signed envelope in version 2 metadata.
type Version2Envelope struct {
	Schema            string           `json:"schema" toml:"schema"`
	LicensorSignature string           `json:"licensorSignature" toml:"licensorSignature"`
	AgentSignature    string           `json:"agentSignature" toml:"agentSignature"`
	Manifest          Version2Manifest `json:"manifest" toml:"manifest"`
}

// Version2Manifest describes a public license and offer.
type Version2Manifest struct {
	// These declarations must appear in this order so as to
	// serialize in the correct order for signature verification.
	LicensorID string `json:"licensorID" toml:"licensorID"`
	OfferID    string `json:"offerID" toml:"offerID"`
	Terms      string `json:"terms" toml:"terms"`
	Version    string `json:"version" toml:"version"`
}
