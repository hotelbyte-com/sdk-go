package protocol

type CheckAvailReq struct {
	RatePkgId string `json:"ratePkgId" required:"true"`
	SessionOption
	TestOption
}
type CheckAvailResp struct {
	Status      CheckAvailStatus `json:"status"`                                // Status indicates whether the room is available or not, or other unavailable causes like credit warning status
	RoomRatePkg *RoomRatePkg     `json:"roomRatePkg,omitzero"`                  // with updated ARI
	Supplier    int64            `json:"supplier,omitempty" apidoc:"HotelCode"` // supplier information for dynamic download
	Header      CommonHeader     `json:"header"`
}
type CheckAvailStatus int

const (
	CheckAvailStatusAvailable   CheckAvailStatus = 1 // CheckAvailStatusAvailable indicates the room is available for booking
	CheckAvailStatusUnavailable CheckAvailStatus = 2 // CheckAvailStatusUnavailable indicates the room is not available
)
