package protocol

import "github.com/hotelbyte-com/sdk-go/protocol/types"

type CancelReq struct {
	CustomerReferenceNo string `json:"customerReferenceNo" required:"true"`
	SupplierReferenceNo string `json:"supplierReferenceNo" required:"true"`
	TestOption
}

type CancelResp struct {
	ServiceFee types.Money `json:"serviceFee"` // ServiceFee is the service fee charged by supplier, not refunded
	Status     OrderStatus `json:"status"`     // Status indicates the current status of the order after cancellation
}
