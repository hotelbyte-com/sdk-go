package protocol

type GuestPerRoom struct {
	AdultCount   int64   `json:"adultCount" required:"true" example:"2"`      // AdultCount is the number of adults
	ChildrenAges []int64 `json:"childrenAges" required:"false" example:"[2]"` // ChildrenAges contains ages of children
	Guests       []Guest `json:"guests" apidoc:"HotelCode"`                   // Guests contains detailed guest information
}

type Guest struct {
	RoomIndex       int64  `json:"roomIndex" required:"true"`                          // Assigned room index for this guest, starts from 1
	FirstName       string `json:"firstName" required:"true" example:"John2"`          // First name of this guest
	LastName        string `json:"lastName" required:"true" example:"Doe"`             // Last name of this guest
	NationalityCode string `json:"nationalityCode" required:"false" example:"US"`      // Nationality code of this guest
	Age             int64  `json:"age,omitempty" required:"false" example:"18"`        // Age of this guest, only matters for children
	IsChild         bool   `json:"isChild,omitempty" required:"false" example:"false"` // Indicates if this guest is a child, determined by clients themselves
}

// Occupancies contains guest and room configuration for hotel searches
type Occupancies struct {
	// CountryCode is the country code of the booker's point of sale in ISO 3166-1 alpha-2 format (e.g., "US")
	CountryCode string `json:"countryCode" required:"true" default:"US" example:"US"`
	// ResidencyCode is the residency code of the booker in ISO 3166-1 alpha-2 format (e.g., "US")
	ResidencyCode string `json:"residencyCode" required:"true" default:"US" example:"US"`
	// NationalityCode is the nationality code of the booker in ISO 3166-1 alpha-2 format (e.g., "US")
	NationalityCode string         `json:"nationalityCode" required:"true" default:"US" example:"US"`
	RoomOccupancies []GuestPerRoom `json:"roomOccupancies"`
}

func (o Occupancies) GetAdultCount() int64 {
	if len(o.RoomOccupancies) <= 0 {
		return 1
	}
	cnt := int64(0)
	for _, room := range o.RoomOccupancies {
		if room.AdultCount > 0 {
			cnt += room.AdultCount
		}
	}
	return cnt
}

func (o Occupancies) GetChildrenCount() int64 {
	if len(o.RoomOccupancies) <= 0 {
		return 0
	}
	cnt := 0
	for _, room := range o.RoomOccupancies {
		cnt += len(room.ChildrenAges)
	}
	return int64(cnt)
}

func (o Occupancies) GetRoomCount() int64 {
	if len(o.RoomOccupancies) <= 0 {
		return 1
	}
	return int64(len(o.RoomOccupancies))
}
