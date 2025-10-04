package protocol

import (
	"github.com/hotelbyte-com/sdk-go/protocol/types"
	"time"
)

type HotelStaticProfile struct {
	LastRefreshStaticTime time.Time `json:"lastRefreshStaticTime,omitzero" apidoc:"HotelCode"` // 最后一次静态信息同步时间
	DestinationId         types.ID  `json:"destinationId,omitzero" example:"804028047"`        // index

	Name              types.I18N              `json:"name,omitzero" example:"{\"en\":\"Jumeirah Beach Hotel\"}"` // hotel name
	Star              float64                 `json:"star,omitzero" example:"4" apidoc:"HotelCode"`              // Star, hotel star, range [3,4,5,6,7]
	LoyaltyProgram    HotelLoyaltyProgram     `json:"loyaltyProgram,omitzero"  apidoc:"HotelCode"`
	Address           types.I18N              `json:"address,omitzero" example:"{\"en\":\"Millenium Al Barsha\"}"`
	LatlngCoordinator types.LatlngCoordinator `json:"latlngCoordinator,omitzero"`
	LogoURL           string                  `json:"logoURL,omitempty" example:"https://static.hotelbyte.com/Images/Hotel/804028047_1.jpg"`
	Phone             string                  `json:"phone,omitempty"  apidoc:"HotelCode"`      // 电话号码
	Email             string                  `json:"email,omitempty"  apidoc:"HotelCode"`      // 电子邮箱
	Fax               string                  `json:"fax,omitempty"  apidoc:"HotelCode"`        // 传真号码
	ZipCode           string                  `json:"zipCode,omitempty"  apidoc:"HotelCode"`    // 邮编
	WebsiteURL        string                  `json:"websiteURL,omitempty"  apidoc:"HotelCode"` // 网站URL

	Desc        types.I18N `json:"desc,omitzero" apidoc:"HotelCode"`
	OpenYear    int64      `json:"openYear,omitempty" apidoc:"HotelCode"`    // 开业年份
	FitmentYear int64      `json:"fitmentYear,omitempty" apidoc:"HotelCode"` // 装修年份

	LastNameUpdateTime time.Time `json:"lastNameUpdateTime,omitzero" apidoc:"HotelCode"` // 最后一次名称变更时间
}
type HotelLoyaltyProgram struct {
	GroupId   int64      `json:"groupId,omitzero"`
	BrandId   int64      `json:"brandId,omitzero"`
	GroupName types.I18N `json:"groupName,omitempty"`
	BrandName types.I18N `json:"brandName,omitempty"`
}
