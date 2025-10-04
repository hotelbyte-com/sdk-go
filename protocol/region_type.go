package protocol

// DestinationType defines the enumeration of destination types
type DestinationType int

const (
	DestinationType_Unknown DestinationType = 0
	// DestinationType_Continent
	//apidoc:zh 大片陆地，例如欧洲、亚洲或南美洲。有关大洲的列表，请参考以下地区概要部分。
	// Continent represents a large landmass, such as Europe, Asia, or South America.
	// For a list of continents, refer to the region summary section below.
	DestinationType_Continent DestinationType = 1
	// DestinationType_Country
	//apidoc:zh 民族国家，例如美国、澳大利亚或德国。国家/地区是从地理意义上定义的，而不一定从政治意义上定义。
	//apidoc:zh 例如，瓜德罗普是加勒比海中的一个法语省。此地区在政治上属于法国，但拥有自己的地理实体。
	//apidoc:zh 有关国家/地区的列表，请参考以下地区概要部分。
	// Country represents a nation-state, such as the United States, Australia, or Germany.
	// Countries/regions are defined geographically, not necessarily politically.
	// For example, Guadeloupe is a French department in the Caribbean.
	// This region is politically part of France but has its own geographical entity.
	// For a list of countries/regions, refer to the region summary section below.
	DestinationType_Country DestinationType = 2
	// DestinationType_ProvinceState
	//apidoc:zh 某个国家内的行政区，例如加利福尼亚州或新南威尔士省。
	//apidoc:zh 包括下级行政区，例如某个国家/地区内的县和省。
	// ProvinceState represents an administrative division within a country,
	// such as California or New South Wales.
	// It includes lower-level administrative divisions, such as counties and provinces within a country/region.
	DestinationType_ProvinceState DestinationType = 3
	// DestinationType_HighLevelRegion
	//apidoc:zh 旅游地区，例如德国的莱茵河上游或者加利福尼亚州的酒乡。
	// HighLevelRegion represents a tourist region, such as the Upper Rhine in Germany
	// or the wine country in California.
	DestinationType_HighLevelRegion DestinationType = 4
	// DestinationType_MultiCityVicinity
	//apidoc:zh 包含一个主要城市及其周边郊区的大城市区域，例如开罗及其周边地区。
	//apidoc:zh 根据旅游价值，它们的多边形边界的定义较为宽泛。
	// MultiCityVicinity represents a large metropolitan area that includes a major city
	// and its surrounding suburbs, such as Cairo and its surrounding areas.
	// Their polygon boundaries are defined broadly based on tourist value.
	DestinationType_MultiCityVicinity DestinationType = 5
	// DestinationType_City
	//apidoc:zh 城市地区，例如贝尔维尤和华盛顿。大城市区域中可能会有一个城市，例如伦敦。
	//apidoc:zh 在此情况下，该城市将覆盖一个比相关 multi_city_vicinity 小的区域。
	//apidoc:zh 与 Expedia 签订合同的每家住宿必须与某个城市关联，因为这是我们地理覆盖的基本单位。
	//apidoc:zh 城市多边形边界是根据官方管理来源定义的。
	// City represents an urban area, such as Bellevue and Washington.
	// A large metropolitan area may contain a city, such as London.
	// In this case, the city will cover a smaller area than the relevant multi_city_vicinity.
	// Each accommodation contracted with Expedia must be associated with a city,
	// as this is the basic unit of our geographical coverage.
	// City polygon boundaries are defined based on official administrative sources.
	DestinationType_City DestinationType = 6
	// DestinationType_Neighborhood
	//apidoc:zh 比城市小的地理区域，例如德国波恩的老城区。
	//apidoc:zh 邻近地区多边形边界是根据旅游价值定义的，而不是管理上的定义。
	// Neighborhood represents a geographical area smaller than a city,
	// such as the old town in Bonn, Germany.
	// Neighborhood polygon boundaries are defined based on tourist value,
	// not administrative definitions.
	DestinationType_Neighborhood DestinationType = 7
	// DestinationType_Airport
	//apidoc:zh 机场，例如 LHR（伦敦希思罗）
	// Airport represents an airport, such as LHR (London Heathrow).
	DestinationType_Airport DestinationType = 8
	// DestinationType_PointOfInterest
	//apidoc:zh 地点，例如博物馆、会议中心、海滩或公园。
	//apidoc:zh 此列表不仅仅包含显而易见的旅游景点，还涵盖范围更广的地点，
	//apidoc:zh 例如上海展览中心、德克萨斯州国会大厦或布拉格动物园。
	// PointOfInterest represents a place, such as a museum, convention center, beach, or park.
	// This list includes not only obvious tourist attractions but also a wider range of places,
	// such as the Shanghai Exhibition Center, the Texas State Capitol, or the Prague Zoo.
	DestinationType_PointOfInterest DestinationType = 9
	// DestinationType_TrainStation
	//apidoc:zh 火车站，例如纽约中央火车站。
	// TrainStation represents a train station, such as Grand Central Terminal in New York.
	DestinationType_TrainStation DestinationType = 10
	// DestinationType_MetroStation
	//apidoc:zh 地铁站，例如俄罗斯圣彼得堡阿威托沃站。
	// MetroStation represents a subway station, such as Avtovo Station in St. Petersburg, Russia.
	DestinationType_MetroStation DestinationType = 11
	// DestinationType_BusStation
	//apidoc:zh 巴士站，例如位于美国纽约州的纽约港务局巴士总站。
	// BusStation represents a bus station, such as the Port Authority Bus Terminal in New York State, USA.
	DestinationType_BusStation DestinationType = 12
)

// String 为 Type 类型实现 String 方法，用于返回枚举值的字符串表示
// String implements the String method for the Type type to return the string representation of the enum value
func (t DestinationType) String() string {
	switch t {
	case DestinationType_Continent:
		return "continent"
	case DestinationType_Country:
		return "country"
	case DestinationType_ProvinceState:
		return "province_state"
	case DestinationType_HighLevelRegion:
		return "high_level_region"
	case DestinationType_MultiCityVicinity:
		return "multi_city_vicinity"
	case DestinationType_City:
		return "city"
	case DestinationType_Neighborhood:
		return "neighborhood"
	case DestinationType_Airport:
		return "airport"
	case DestinationType_PointOfInterest:
		return "point_of_interest"
	case DestinationType_TrainStation:
		return "train_station"
	case DestinationType_MetroStation:
		return "metro_station"
	case DestinationType_BusStation:
		return "bus_station"
	default:
		return "unknown"
	}
}
