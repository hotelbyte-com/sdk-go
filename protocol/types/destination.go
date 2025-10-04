package types

type LatlngCoordinator struct {
	Google *Latlng `json:"google,omitempty" required:"true" example:"{\"lat\":25.0478,\"lng\":121.5319}"`
	Gaode  *Latlng `json:"gaode,omitempty" example:"{\"lat\":25.0478,\"lng\":121.5319}"`
}
type Latlng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
