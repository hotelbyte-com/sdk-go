package protocol

import (
	"encoding/json"
	"testing"

	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

func TestDistributorOptionSerializesEntityIDAsString(t *testing.T) {
	raw, err := json.Marshal(HotelListReq{DistributorOption: DistributorOption{DistributorID: types.ID(42)}})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) == "" || !containsJSON(raw, `"distributorId":"42"`) {
		t.Fatalf("unexpected distributor JSON: %s", raw)
	}
}

func containsJSON(raw []byte, want string) bool {
	for i := 0; i+len(want) <= len(raw); i++ {
		if string(raw[i:i+len(want)]) == want {
			return true
		}
	}
	return false
}
