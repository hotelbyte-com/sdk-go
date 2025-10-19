package protocol

import (
	"time"

	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

type HotelListReq struct {
	HotelIds         types.IDs `json:"hotelIds" example:"[\"461850557\",\"118062388\"]"`
	MaxRatesPerHotel int64     `json:"maxRatesPerHotel" default:"0" example:"3"`
	CheckInOut                 // embedded
	Occupancies                // embedded
	HotelDestination           // embedded
	HotelListFilter            // embedded
	types.PageReq              // embedded
	SortBy           string    `json:"sortBy,omitempty" example:"price-desc"`

	// header
	CurrencyOption
	TestOption
}

type HotelListResp struct {
	List  HotelList          `json:"list,omitzero" required:"true"`
	Basic HotelListBasicInfo `json:"basic,omitzero" required:"true"`
	types.PageResp
}

type HotelList []*Hotel

type Hotel struct {
	ID types.ID `json:"id,omitzero" example:"461850557"`
	HotelStaticProfile
	MinPrice    types.Money `json:"minPrice,omitzero" example:"{\"amount\":100,\"currency\":\"USD\"}"` // the minPrice meeting search criteria
	IsAvailable bool        `json:"isAvailable,omitempty" example:"true"`
	Rooms       []Room      `json:"rooms,omitzero"`
}
type Room struct {
	RoomTypeId   string        `json:"roomTypeId,omitempty" required:"true"`  // standardized roomTypeId, e.g. "R001"
	RoomTypeName types.I18N    `json:"roomTypeName,omitzero" required:"true"` // standardized room type name; if not recognized, we'll give a default name like "Standard Room"
	HotelId      types.ID      `json:"hotelId,omitempty" required:"true"`     // hotelId, duplicate of parent HotelId for indicating identified room
	Rates        []RoomRatePkg `json:"rates,omitzero"`                        // core structure, room offers under this room
}

type RoomRatePkg struct {
	RatePkgId string `json:"ratePkgId,omitempty" required:"true"` // used as key input for checkAvail & book APIs
	ComputedCancelPolicy
	OriginalRoomNaming OriginalRoomNaming `json:"originalRoomNaming,omitzero"`   // origin room naming fields from supplier
	Rate               Rate               `json:"rate,omitzero" required:"true"` // single room price for multiple nights
	TotalRate          Rate               `json:"totalRate,omitzero"`            // total price for multiple rooms (and multiple nights)
	RateComment        string             `json:"rateComment"`
	RatePlan
}

type ComputedCancelPolicy struct {
	RefundableMode  RefundableMode             `json:"refundableMode,omitempty" required:"true"`
	RefundableUntil time.Time                  `json:"refundableUntil,omitzero"` //todo
	CancelFees      []ComputedCancelPolicyItem `json:"cancelFees,omitzero"`
}
type OriginalRoomNaming struct {
	Id       string `json:"id,omitempty"`       // unique identifier (or say, room code)
	Name     string `json:"name,omitempty"`     // room name comes from supplier
	Supplier int64  `json:"supplier,omitempty"` // supplier id
}

type ComputedCancelPolicyItem struct {
	Until time.Time   `json:"until,omitzero"` // cancellation deadline
	Fee   types.Money `json:"fee,omitzero"`   // cancellation fee, 0 means free cancellation
}

type RefundableMode string

const (
	RefundableModeFully     RefundableMode = "full"    // free cancellation
	RefundableModePartially RefundableMode = "partial" // partial cancellation
	RefundableModeNo        RefundableMode = "no"      // not refundable
)

func (r RefundableMode) Bool() bool {
	return r != RefundableModeNo
}

type RatePlan struct {
	// Taxes and fees, generally collected by hotels on behalf of the government, are fixed costs and do not participate in the price increase during the distribution process. Therefore, they are listed separately.
	Tax Tax `json:"tax,omitzero"`
}
type TaxItem struct {
	TaxType string      `json:"taxType"` // VAT, SERVICE_CHARGE, CITY_TAX, etc.
	TaxName string      `json:"taxName"`
	Amount  types.Money `json:"amount"`
	Desc    string      `json:"desc,omitempty"`
}
type Tax struct {
	Total types.Money `json:"total,omitzero"`
	Items []TaxItem   `json:"items,omitzero"`
}
type HotelListBasicInfo struct {
	DestinationId int64  `json:"destinationId,omitempty"`
	SessionId     string `json:"sessionId,omitempty"`
}
