package protocol

import "github.com/hotelbyte-com/sdk-go/protocol/types"

type HotelRatesReq struct {
	HotelId types.ID `json:"hotelId" required:"true" example:"461850557"`
	CheckInOut
	Occupancies
	HotelDestination
	CurrencyOption
	SessionOption
	TestOption
}

type CurrencyOption struct {
	// Requested currency for the rates, in https://en.wikipedia.org/wiki/ISO_4217 format
	Currency string `json:"currency,omitempty" api.header:"Currency"`
}
type SessionOption struct {
	// Suggested provided in request by client. It's required in booking flow.
	SessionId string `json:"sessionId,omitempty" api.header:"Session-Id" required:"true"`
}
type TestOption struct {
	Test string `json:"test" api.header:"Test"` // Test flags. support key-value pairs, eg, "hotel=HC1&scenario=priceChange".If it's not recognized by server, the call will behave as if the "Test" header was not provided.
}
type HotelRatesResp struct {
	Rooms  []*Room      `json:"rooms"`
	Header CommonHeader `json:"header"`
}

type CommonHeader struct {
	RequestId string `json:"requestId,omitempty" api.header:"Request-Id"` // for identifying the current request, it can't be duplicate
	TraceId   string `json:"traceId,omitempty" api.header:"Trace-Id"`     // for tracing a group of requests
}
