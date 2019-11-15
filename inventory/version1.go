package inventory

// Version1Envelope describes a signed manifest in licensezero.json and Cargo.toml files.
type Version1Envelope struct {
	LicensorSignature string           `json:"licensorSignature" toml:"licensorSignature"`
	AgentSignature    string           `json:"agentSignature" toml:"agentSignature"`
	Manifest          Version1Manifest `json:"license" toml:"license"`
}

// Version1Manifest describes a public license and offer.
type Version1Manifest struct {
	// Note: These declarations must appear in this order so as to
	// serialize in the correct order for signature verification.
	Repository   string `json:"homepage" toml:"homepage"`
	Jurisdiction string `json:"jurisdiction" tom:"jurisdiction"`
	Name         string `json:"name" toml:"name"`
	ProjectID    string `json:"offerID" toml:"offerID"`
	PublicKey    string `json:"publicKey" toml:"publicKey"`
	Terms        string `json:"terms" toml:"terms"`
	Version      string `json:"version" toml:"version"`
}
