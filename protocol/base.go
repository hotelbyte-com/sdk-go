package protocol

import (
	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

type CheckInOut struct {
	CheckIn  types.DateInt `json:"checkIn" required:"true" example:"2026-01-01" default:"now()"`    // CheckIn is the check-in date, defaults to current date
	CheckOut types.DateInt `json:"checkOut" required:"true" example:"2026-01-03" default:"now()+7"` // CheckOut is the check-out date, defaults to 7 days from now
}

type HotelDestination struct {
	DestinationName string `json:"destinationName"` // optional
}

type HotelListFilter struct {
	Distance *DistanceFilter `json:"distance,omitempty" apidoc:"HotelCode"` // Distance contains distance-based filtering criteria
	Price    *PriceFilter    `json:"price,omitempty" apidoc:"HotelCode"`    // Price contains price range filtering criteria
}

// DistanceFilter contains distance-based filtering criteria
type DistanceFilter struct {
	Latlng types.LatlngCoordinator `json:"latlng"`
	Radius float64                 `json:"radius" example:"1000"`
}

type PriceFilter struct {
	LowPrice  float64 `json:"lowPrice" example:"100"`   // LowPrice is the minimum price threshold
	HighPrice float64 `json:"highPrice" example:"1000"` // HighPrice is the maximum price threshold
}

type AuthReq struct {
	AppKey    string `json:"appKey"`    // Application key for authentication
	AppSecret string `json:"appSecret"` // Application secret for authentication
	TTL       int64  `json:"ttl"`       // Time-to-live in seconds for the token
}

type AuthResp struct {
	Ticket string `json:"ticket"`
}
