package protocol

import "github.com/hotelbyte-com/sdk-go/protocol/types"

// DistributorOption selects a connected marketplace distributor. Search
// requests omit it to aggregate; CheckAvail/Book use the session-bound route
// when omitted, while order operations use the persisted order route.
type DistributorOption struct {
	DistributorID types.ID `json:"distributorId,omitempty"`
}
