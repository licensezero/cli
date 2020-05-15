package api

// Pricing describes private license pricing data.
type Pricing struct {
	Private   uint `json:"private"`
	Relicense uint `json:"relicense,omitempty"`
}
