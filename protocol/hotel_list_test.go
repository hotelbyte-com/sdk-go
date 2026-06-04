package protocol

import (
	"testing"

	"github.com/bytedance/sonic"
)

func TestMarshal(t *testing.T) {
	req := &HotelListReq{}
	s, _ := sonic.MarshalString(req)
	t.Logf("%+v", s)
}
