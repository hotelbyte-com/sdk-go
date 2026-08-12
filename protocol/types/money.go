package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type Money struct {
	Currency string  `json:"currency" required:"true"` // eg. "USD"
	Amount   float64 `json:"amount" required:"true"`   // eg. "14.50"
}

// MarshalJSON emits the canonical decimal-string amount used by hotel-be.
// The public Amount field remains float64 for SDK source compatibility.
func (m Money) MarshalJSON() ([]byte, error) {
	if math.IsNaN(m.Amount) || math.IsInf(m.Amount, 0) {
		return nil, fmt.Errorf("Money.Amount: cannot marshal non-finite value %v", m.Amount)
	}
	return json.Marshal(struct {
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	}{
		Currency: m.Currency,
		Amount:   strconv.FormatFloat(m.Amount, 'f', -1, 64),
	})
}

// UnmarshalJSON accepts the canonical decimal-string amount emitted by
// hotel-be and the legacy JSON number representation. Amount remains float64
// for public API compatibility.
func (m *Money) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("Money: nil receiver")
	}
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*m = Money{}
		return nil
	}

	var wire struct {
		Currency string          `json:"currency"`
		Amount   json.RawMessage `json:"amount"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	amount, err := parseMoneyAmount(wire.Amount)
	if err != nil {
		return err
	}
	*m = Money{Currency: wire.Currency, Amount: amount}
	return nil
}

func parseMoneyAmount(raw json.RawMessage) (float64, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return 0, nil
	}

	var decimalString string
	if err := json.Unmarshal(trimmed, &decimalString); err == nil {
		amount, parseErr := strconv.ParseFloat(decimalString, 64)
		if parseErr != nil {
			return 0, fmt.Errorf("Money.Amount: invalid decimal string %q: %w", decimalString, parseErr)
		}
		if math.IsNaN(amount) || math.IsInf(amount, 0) {
			return 0, fmt.Errorf("Money.Amount: non-finite decimal string %q", decimalString)
		}
		return amount, nil
	}

	var number float64
	if err := json.Unmarshal(trimmed, &number); err == nil {
		return number, nil
	}
	return 0, fmt.Errorf("Money.Amount: expected string or number, got %s", trimmed)
}
