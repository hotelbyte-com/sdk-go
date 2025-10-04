package protocol

import (
	"github.com/hotelbyte-com/sdk-go/protocol/types"
	"time"
)

type BookReq struct {
	CustomerReferenceNo string  `json:"customerReferenceNo,omitempty" example:"uuid"` // customerReferenceNo contains the reference number value
	RatePkgId           string  `json:"ratePkgId,omitempty" required:"true"`          // RatePkgId is obtained from HotelStaticDetail API
	Holder              Holder  `json:"holder,omitzero" required:"true"`              // Holder contains the booking contact information
	Guests              []Guest `json:"guests,omitzero" required:"true"`              // Guests contains the list of guests for this room
	SessionOption
}
type Holder struct {
	FirstName string `json:"firstName" required:"true" example:"John"`
	LastName  string `json:"lastName" required:"true" example:"Doe"`
	Email     string `json:"email,omitempty" required:"false" example:"John@hotelbyte.com"`
	Phone     Phone  `json:"phone,omitzero"`
}
type Phone struct {
	CountryCode   string `json:"countryCode,omitempty"`   // AE
	CountryNumber int64  `json:"countryNumber,omitempty"` // 971
	Number        string `json:"number,omitempty"`        // 525757249
}
type BookResp struct {
	HotelOrder *HotelOrder `json:"hotelOrder,omitzero"` // HotelOrder contains the hotel order information
}

// OrderHotelInfo contains hotel information associated with an order
type OrderHotelInfo struct {
	HotelId types.ID `json:"hotelId"`
	HotelStaticProfile
}

// OrderRoomInfo contains detailed room information for an order
type OrderRoomInfo struct {
	Room
	RoomRatePkg
	RoomIndex  int64                 `json:"roomIndex"` // room index
	Guests     []Guest               `json:"guests"`    // information about all guests in this room
	RefundInfo []OrderRoomRefundInfo `json:"refundInfo,omitzero"`
}

// OrderRoomRefundInfo represents refund information for a specific room on a specific date
// It tracks whether a refund has been processed and the amount refunded
type OrderRoomRefundInfo struct {
	Date          types.DateInt `json:"date" required:"true"`          // Date for which the refund applies
	Refunded      bool          `json:"refunded" required:"true"`      // Whether the refund has been processed
	RefundedMoney types.Money   `json:"refundedMoney" required:"true"` // Amount of money refunded for this date
}

// HotelOrder represents a complete hotel order with all associated information
type HotelOrder struct {
	*OrderBasic
	Hotel *OrderHotelInfo  `json:"hotel,omitempty"` // hotel-specific information
	Rooms []*OrderRoomInfo `json:"rooms"`           // detailed information for each booked room
}

// OrderBasic contains the fundamental information about a hotel order
type OrderBasic struct {
	Status              OrderStatus   `json:"status" required:"true" example:"1"` // OrderStatus indicates the current status of the order
	CheckIn             types.DateInt `json:"checkIn" required:"true" example:"2026-01-01"`
	CheckOut            types.DateInt `json:"checkOut" required:"true" example:"2026-01-03"`
	NightCount          int64         `json:"nightCount,omitzero"`                                                                                                                                                           // Number of nights for the stay
	RoomCount           int64         `json:"roomCount,omitzero"`                                                                                                                                                            // CheckOut is the check-out date for the booking
	BookingTime         time.Time     `json:"bookingTime" required:"true" example:"2026-01-01T10:00:00Z"`                                                                                                                    // BookingTime is when the booking was created
	HotelConfirmNo      string        `json:"hotelConfirmNo" required:"false" example:"123456"`                                                                                                                              // HotelConfirmNo is the confirmation number from the hotel
	Holder              Holder        `json:"holder" required:"true" example:"{\"firstName\":\"John\",\"lastName\": \"Doe\",\"phone\":{\"countryCode\":\"1\",\"number\":\"1234567890\"},\"email\":\"johndoe@example.com\"}"` // Holder contains the person who made the booking
	CustomerReferenceNo string        `json:"customerReferenceNo" required:"false" example:"1234567890"`                                                                                                                     // CustomerReferenceNo is the unique order identifier from the customer
	CancelTime          time.Time     `json:"cancelTime,omitzero" example:"2026-01-02T10:00:00Z"`                                                                                                                            // CancelTime is when the booking was cancelled, if applicable
	CancelReason        string        `json:"cancelReason" example:"Customer requested cancellation"`                                                                                                                        // CancelReason contains the reason if the order was cancelled
	RefundedPrice       types.Money   `json:"refundedPrice,omitzero" example:"{\"amount\":99.00,\"currency\":\"USD\"}"`                                                                                                      // RefundedPrice is the amount that has been refunded to the customer
	Supplier            int64         `json:"supplier" example:"20000001"`
	SupplierReferenceNo string        `json:"supplierReferenceNo" required:"true"` // SupplierReferenceNo is the unique order identifier from the supplier	// Supplier identifies which supplier processed this order
	Rate
}

type Rate struct {
	CommissionableRate types.Money `json:"commissionableRate,omitzero" required:"false"` // used for commission calculation
	NetRate            types.Money `json:"netRate,omitzero" required:"true"`
	GrossRate          types.Money `json:"grossRate,omitzero" required:"false"`
	RespectGrossRate   bool        `json:"respectGrossRate,omitzero" required:"false"` // You should respect GrossRate if RespectGrossRate is true; default as false
}

type OrderStatus int

const (
	OrderStatus_Unknown OrderStatus = iota
	OrderStatus_Confirming
	OrderStatus_Confirmed
	OrderStatus_Cancelled
	OrderStatus_Failed
	OrderStatus_CancelFailed
)
