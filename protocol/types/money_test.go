package types

import (
	"encoding/json"
	"math"
	"testing"
)

func TestMoneyMarshalJSONEmitsCanonicalStringAmount(t *testing.T) {
	raw, err := json.Marshal(Money{Currency: "USD", Amount: 198.6})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"currency":"USD","amount":"198.6"}` {
		t.Fatalf("money JSON = %s", raw)
	}
}

func TestMoneyJSONRoundTrip(t *testing.T) {
	want := Money{Currency: "AED", Amount: 14.5}
	raw, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got Money
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("money round trip = %+v, want %+v", got, want)
	}
}

func TestMoneyMarshalJSONRejectsNonFiniteAmount(t *testing.T) {
	for _, amount := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := json.Marshal(Money{Currency: "USD", Amount: amount}); err == nil {
			t.Fatalf("expected amount %v to fail", amount)
		}
	}
}

func TestMoneyUnmarshalJSONAcceptsCanonicalStringAmount(t *testing.T) {
	var money Money
	if err := json.Unmarshal([]byte(`{"currency":"USD","amount":"198.6"}`), &money); err != nil {
		t.Fatal(err)
	}
	if money.Currency != "USD" || money.Amount != 198.6 {
		t.Fatalf("money = %+v", money)
	}
}

func TestMoneyUnmarshalJSONAcceptsLegacyNumberAmount(t *testing.T) {
	var money Money
	if err := json.Unmarshal([]byte(`{"currency":"AED","amount":14.5}`), &money); err != nil {
		t.Fatal(err)
	}
	if money.Currency != "AED" || money.Amount != 14.5 {
		t.Fatalf("money = %+v", money)
	}
}

func TestMoneyUnmarshalJSONRejectsInvalidAmount(t *testing.T) {
	for _, raw := range []string{
		`{"currency":"USD","amount":"not-a-decimal"}`,
		`{"currency":"USD","amount":"NaN"}`,
		`{"currency":"USD","amount":"Inf"}`,
		`{"currency":"USD","amount":{}}`,
	} {
		var money Money
		if err := json.Unmarshal([]byte(raw), &money); err == nil {
			t.Fatalf("expected invalid amount %s to fail", raw)
		}
	}
}
