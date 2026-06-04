package protocol

import (
	"encoding/json"
	"testing"
)

func TestMarshal(t *testing.T) {
	req := &HotelListReq{}
	b, _ := json.Marshal(req)
	t.Logf("%+v", string(b))
}
