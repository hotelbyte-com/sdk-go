package protocol

import "github.com/hotelbyte-com/sdk-go/protocol/types"

type HotelRatesReq struct {
	HotelId types.ID `json:"hotelId" required:"true" example:"461850557"`
	CheckInOut
	Occupancies
	HotelDestination
	CurrencyOption
	SessionOption
}

type CurrencyOption struct {
	// Requested currency for the rates, in https://en.wikipedia.org/wiki/ISO_4217 format
	Currency string `json:"currency,omitempty" api.header:"Currency"`
}
type SessionOption struct {
	// Suggested provided in request by client. It's required in booking flow.
	SessionId string `json:"sessionId,omitempty" api.header:"Session-Id" required:"true"`
}
type HotelRatesResp struct {
	Rooms []*Room `json:"rooms"`
}
