package inventory

// Version1Envelope describes a signed manifest in licensezero.json and Cargo.toml files.
type Version1Envelope struct {
	LicensorSignature string           `json:"licensorSignature" toml:"licensorSignature" mapstructure:"licensorSignature"`
	AgentSignature    string           `json:"agentSignature" toml:"agentSignature" mapstructure:"agentSignature"`
	Manifest          Version1Manifest `json:"license" toml:"license" mapstructure:"license"`
}

// Version1Manifest describes a public license and offer.
type Version1Manifest struct {
	// Note: These declarations must appear in this order so as to
	// serialize in the correct order for signature verification.
	Repository   string `json:"homepage" toml:"homepage" mapstructure:"homepage"`
	Jurisdiction string `json:"jurisdiction" tom:"jurisdiction" mapstructure:"jurisdiction"`
	Name         string `json:"name" toml:"name" mapstructure:"name"`
	ProjectID    string `json:"projectID" toml:"projectID" mapstructure:"projectID"`
	PublicKey    string `json:"publicKey" toml:"publicKey" mapstructure:"publicKey"`
	Terms        string `json:"terms" toml:"terms" mapstructure:"terms"`
	Version      string `json:"version" toml:"version" mapstructure:"version"`
}

func (envelope *Version1Envelope) offer() Offer {
	var manifest = envelope.Manifest
	return Offer{
		License: LicenseData{
			Terms:   manifest.Terms,
			Version: manifest.Version,
		},
		OfferID:  manifest.ProjectID,
		Envelope: envelope,
	}
}
