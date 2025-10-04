package protocol

import "github.com/hotelbyte-com/sdk-go/protocol/types"

// QueryOrdersReq represents a request to query multiple hotel orders
type QueryOrdersReq struct {
	CustomerReferenceNos []string          `json:"customerReferenceNos,omitzero"` // customer reference numbers to search for
	SupplierReferenceNos []string          `json:"supplierReferenceNos,omitzero"` // supplier reference numbers to search for
	CheckInTimeWindow    *types.TimeWindow `json:"checkInTimeWindow,omitzero"`    // filters orders by check-in date range
	CheckOutTimeWindow   *types.TimeWindow `json:"checkOutTimeWindow,omitzero"`   // filters orders by check-out date range
	BookingTimeWindow    *types.TimeWindow `json:"bookingTimeWindow,omitzero"`    // filters orders by creation date range
	FreeCancelTimeWindow *types.TimeWindow `json:"freeCancelTimeWindow,omitzero"` // filters orders by free cancel date range
	CancelledTimeWindow  *types.TimeWindow `json:"cancelledTimeWindow,omitzero"`  // filters orders by cancellation date range
	StatusList           []OrderStatus     `json:"statusList,omitzero"`           // filters orders by status
}

// QueryOrdersResp represents the response containing multiple hotel orders
type QueryOrdersResp struct {
	Orders []*HotelOrder `json:"orders"` // Orders contains a list of hotel order information
}
