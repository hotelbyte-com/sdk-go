package types

type Money struct {
	Currency string  `json:"currency" required:"true"` // eg. "USD"
	Amount   float64 `json:"amount" required:"true"`   // eg. "14.50"
}
