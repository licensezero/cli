package api

// Price represents the price paid for the license.
type Price struct {
	// Fields of this struct must be sorted by JSON key, in
	// order to serialize correctly for signature.
	Amount   uint   `json:"amount"`
	Currency string `json:"currency"`
}
