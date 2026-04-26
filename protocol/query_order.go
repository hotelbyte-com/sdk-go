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
	TestOption
}

// QueryOrdersResp represents the response containing multiple hotel orders.
//
// MVP order-detail field paths:
// - orders[].rooms[0].includesPackaging indicates packaged rates.
// - orders[].rooms[0].totalRate.grossRate is gross price; display it only when respectGrossRate is true and amount is non-zero.
// - orders[].rooms[0].rateComment contains rate comments.
// - orders[].sourceRate.netRate is supplier original net rate when exposed by the API.
// - orders[].rooms[0].totalRateInBusinessCurrency.netRate is customer net in business currency when exposed by the API.
// - orders[].rooms[0].policyBufferHours and orders[].rooms[0].cancellationPolicyText are backend cancellation display fields.
// - channel/distribution should be derived as: bookingSource XML => API; Retail/Admin => Web; bookingMode external => Offline, otherwise Direct.
type QueryOrdersResp struct {
	Orders []*HotelOrder `json:"orders"` // Orders contains a list of hotel order information
	Header CommonHeader  `json:"header"`
}
