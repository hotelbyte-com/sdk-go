package protocol

import "github.com/hotelbyte-com/sdk-go/protocol/types"

type Destination struct {
	ID                     types.ID        `json:"id"`
	Type                   DestinationType `json:"type,omitempty"`
	Name                   types.I18N      `json:"name,omitempty"`                                      // i18n name
	NameFull               types.I18N      `json:"nameFull,omitempty"`                                  // eg."Springfield, Missouri, United States of America"
	CountryCode            string          `json:"countryCode,omitempty"`                               // specified ISO 3166-1 alpha-2 country code, see https://www.iso.org/obp/ui/#search
	CountrySubdivisionCode string          `json:"countrySubdivisionCode,omitempty" apidoc:"HotelCode"` // ISO 3166-2 country subdivision

	// transient field
	Ancestors             []*Destination `json:"ancestors,omitzero" apidoc:"HotelCode"`   // An array of the region's ancestors. Sorted from root.
	Descendants           []*Destination `json:"descendants,omitzero" apidoc:"HotelCode"` // Optional. An array of the region's descendants.
	CountryName           types.I18N     `json:"countryName,omitempty"`
	DestinationName       types.I18N     `json:"destinationName,omitempty"`
	ParentDestinationName types.I18N     `json:"parentDestinationName,omitempty"`
	ParentDestinationId   string         `json:"parentDestinationId,omitempty"` // Parent region code, e.g., for city, it's province code; for province, it's country code
}
