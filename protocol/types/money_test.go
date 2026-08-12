package types

import (
	"encoding/json"
	"testing"
)

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
