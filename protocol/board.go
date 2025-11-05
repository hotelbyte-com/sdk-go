package protocol

// BoardId represents the meal plan type following liteapi standard
// See: https://docs.liteapi.travel/docs/hotel-rates-api-json-data-structure
type BoardId string

const (
	// ============================================================
	// 基础餐食类型 - Basic Meal Plan Types
	// ============================================================
	BoardIdRoomOnly     BoardId = "RO" // Room Only - 仅住宿
	BoardIdBedBreakfast BoardId = "BB" // Bed and Breakfast - 含早餐
	BoardIdHalfBoard    BoardId = "HB" // Half Board - 半食宿（早+晚）
	BoardIdFullBoard    BoardId = "FB" // Full Board - 全食宿（三餐）
	BoardIdAllInclusive BoardId = "AI" // All Inclusive - 全包

	// ============================================================
	// 单餐类型 - Single Meal Types
	// ============================================================
	BoardIdBreakfastIncluded BoardId = "BI" // Breakfast Included - 仅早餐
	BoardIdLunchOnly         BoardId = "LU" // Lunch Only - 仅午餐
	BoardIdDinnerOnly        BoardId = "DI" // Dinner Only - 仅晚餐

	// ============================================================
	// 组合餐类型 - Combination Meal Types
	// ============================================================
	BoardIdBreakfastDinner BoardId = "BD" // Breakfast and Dinner - 早+晚
	BoardIdBreakfastLunch  BoardId = "BL" // Breakfast and Lunch - 早+午
	BoardIdLunchDinner     BoardId = "LD" // Lunch and Dinner - 午+晚

	// ============================================================
	// 备用/扩展类型 - Alternate Codes & Extended Types
	// ============================================================
	BoardIdAllInclusiveTI BoardId = "TI" // All Inclusive (alternate code) - 全包（备用代码）

	// ============================================================
	// 特定人数的餐食类型 - Meal Plans for Specific Occupancy
	// ============================================================
	// 1 Person variants
	BoardIdBedBreakfast1 BoardId = "BB1" // Bed and Breakfast for 1 person
	BoardIdHalfBoard1    BoardId = "HB1" // Half Board for 1 person
	BoardIdFullBoard1    BoardId = "FB1" // Full Board for 1 person
	BoardIdAllInclusive1 BoardId = "AI1" // All Inclusive for 1 person

	// 2 Person variants
	BoardIdBedBreakfast2 BoardId = "BB2" // Bed and Breakfast for 2 persons
	BoardIdHalfBoard2    BoardId = "HB2" // Half Board for 2 persons
	BoardIdFullBoard2    BoardId = "FB2" // Full Board for 2 persons
	BoardIdAllInclusive2 BoardId = "AI2" // All Inclusive for 2 persons

	// 3 Person variants
	BoardIdBedBreakfast3 BoardId = "BB3" // Bed and Breakfast for 3 persons
	BoardIdHalfBoard3    BoardId = "HB3" // Half Board for 3 persons
	BoardIdFullBoard3    BoardId = "FB3" // Full Board for 3 persons
	BoardIdAllInclusive3 BoardId = "AI3" // All Inclusive for 3 persons
)

// String returns the string representation of BoardId
func (b BoardId) String() string {
	return string(b)
}

// GetDescription returns a human-readable description of the board type
// Returns description in English and Chinese
func (b BoardId) GetDescription() string {
	switch b {
	// Basic types
	case BoardIdRoomOnly:
		return "Room Only (仅住宿)"
	case BoardIdBedBreakfast, BoardIdBedBreakfast1, BoardIdBedBreakfast2, BoardIdBedBreakfast3:
		return "Bed and Breakfast (含早餐)"
	case BoardIdHalfBoard, BoardIdHalfBoard1, BoardIdHalfBoard2, BoardIdHalfBoard3:
		return "Half Board (半食宿)"
	case BoardIdFullBoard, BoardIdFullBoard1, BoardIdFullBoard2, BoardIdFullBoard3:
		return "Full Board (全食宿)"
	case BoardIdAllInclusive, BoardIdAllInclusiveTI, BoardIdAllInclusive1, BoardIdAllInclusive2, BoardIdAllInclusive3:
		return "All Inclusive (全包)"

	// Single meal types
	case BoardIdBreakfastIncluded:
		return "Breakfast Included (仅早餐)"
	case BoardIdLunchOnly:
		return "Lunch Only (仅午餐)"
	case BoardIdDinnerOnly:
		return "Dinner Only (仅晚餐)"

	// Combination types
	case BoardIdBreakfastDinner:
		return "Breakfast and Dinner (早餐和晚餐)"
	case BoardIdBreakfastLunch:
		return "Breakfast and Lunch (早餐和午餐)"
	case BoardIdLunchDinner:
		return "Lunch and Dinner (午餐和晚餐)"

	default:
		return "Unknown (" + string(b) + ")"
	}
}

// GetNameEn returns the English name of the board type
func (b BoardId) GetNameEn() string {
	switch b {
	case BoardIdRoomOnly:
		return "Room Only"
	case BoardIdBedBreakfast:
		return "Bed and Breakfast"
	case BoardIdBedBreakfast1:
		return "Bed and Breakfast for 1"
	case BoardIdBedBreakfast2:
		return "Bed and Breakfast for 2"
	case BoardIdBedBreakfast3:
		return "Bed and Breakfast for 3"
	case BoardIdHalfBoard:
		return "Half Board"
	case BoardIdHalfBoard1:
		return "Half Board for 1"
	case BoardIdHalfBoard2:
		return "Half Board for 2"
	case BoardIdHalfBoard3:
		return "Half Board for 3"
	case BoardIdFullBoard:
		return "Full Board"
	case BoardIdFullBoard1:
		return "Full Board for 1"
	case BoardIdFullBoard2:
		return "Full Board for 2"
	case BoardIdFullBoard3:
		return "Full Board for 3"
	case BoardIdAllInclusive, BoardIdAllInclusiveTI:
		return "All Inclusive"
	case BoardIdAllInclusive1:
		return "All Inclusive for 1"
	case BoardIdAllInclusive2:
		return "All Inclusive for 2"
	case BoardIdAllInclusive3:
		return "All Inclusive for 3"
	case BoardIdBreakfastIncluded:
		return "Breakfast Included"
	case BoardIdLunchOnly:
		return "Lunch Only"
	case BoardIdDinnerOnly:
		return "Dinner Only"
	case BoardIdBreakfastDinner:
		return "Breakfast and Dinner"
	case BoardIdBreakfastLunch:
		return "Breakfast and Lunch"
	case BoardIdLunchDinner:
		return "Lunch and Dinner"
	default:
		return "Unknown"
	}
}

// GetNameZh returns the Chinese name of the board type
func (b BoardId) GetNameZh() string {
	switch b {
	case BoardIdRoomOnly:
		return "仅住宿"
	case BoardIdBedBreakfast:
		return "含早餐"
	case BoardIdBedBreakfast1:
		return "含早餐 (1人)"
	case BoardIdBedBreakfast2:
		return "含早餐 (2人)"
	case BoardIdBedBreakfast3:
		return "含早餐 (3人)"
	case BoardIdHalfBoard:
		return "半食宿"
	case BoardIdHalfBoard1:
		return "半食宿 (1人)"
	case BoardIdHalfBoard2:
		return "半食宿 (2人)"
	case BoardIdHalfBoard3:
		return "半食宿 (3人)"
	case BoardIdFullBoard:
		return "全食宿"
	case BoardIdFullBoard1:
		return "全食宿 (1人)"
	case BoardIdFullBoard2:
		return "全食宿 (2人)"
	case BoardIdFullBoard3:
		return "全食宿 (3人)"
	case BoardIdAllInclusive, BoardIdAllInclusiveTI:
		return "全包"
	case BoardIdAllInclusive1:
		return "全包 (1人)"
	case BoardIdAllInclusive2:
		return "全包 (2人)"
	case BoardIdAllInclusive3:
		return "全包 (3人)"
	case BoardIdBreakfastIncluded:
		return "仅早餐"
	case BoardIdLunchOnly:
		return "仅午餐"
	case BoardIdDinnerOnly:
		return "仅晚餐"
	case BoardIdBreakfastDinner:
		return "早餐和晚餐"
	case BoardIdBreakfastLunch:
		return "早餐和午餐"
	case BoardIdLunchDinner:
		return "午餐和晚餐"
	default:
		return "未知"
	}
}

// IsValid returns true if the BoardId is a recognized enum value
func (b BoardId) IsValid() bool {
	switch b {
	case
		BoardIdRoomOnly,
		BoardIdBedBreakfast, BoardIdBedBreakfast1, BoardIdBedBreakfast2, BoardIdBedBreakfast3,
		BoardIdHalfBoard, BoardIdHalfBoard1, BoardIdHalfBoard2, BoardIdHalfBoard3,
		BoardIdFullBoard, BoardIdFullBoard1, BoardIdFullBoard2, BoardIdFullBoard3,
		BoardIdAllInclusive, BoardIdAllInclusiveTI, BoardIdAllInclusive1, BoardIdAllInclusive2, BoardIdAllInclusive3,
		BoardIdBreakfastIncluded,
		BoardIdLunchOnly,
		BoardIdDinnerOnly,
		BoardIdBreakfastDinner,
		BoardIdBreakfastLunch,
		BoardIdLunchDinner:
		return true
	default:
		return false
	}
}

// AllBoardIds returns a slice of all defined BoardId enum values
func AllBoardIds() []BoardId {
	return []BoardId{
		// Basic types
		BoardIdRoomOnly,
		BoardIdBedBreakfast,
		BoardIdHalfBoard,
		BoardIdFullBoard,
		BoardIdAllInclusive,

		// Single meal types
		BoardIdBreakfastIncluded,
		BoardIdLunchOnly,
		BoardIdDinnerOnly,

		// Combination types
		BoardIdBreakfastDinner,
		BoardIdBreakfastLunch,
		BoardIdLunchDinner,

		// Alternate codes
		BoardIdAllInclusiveTI,

		// Occupancy variants
		BoardIdBedBreakfast1, BoardIdBedBreakfast2, BoardIdBedBreakfast3,
		BoardIdHalfBoard1, BoardIdHalfBoard2, BoardIdHalfBoard3,
		BoardIdFullBoard1, BoardIdFullBoard2, BoardIdFullBoard3,
		BoardIdAllInclusive1, BoardIdAllInclusive2, BoardIdAllInclusive3,
	}
}

// BoardIdMap returns a map of BoardId code to friendly name (English)
func BoardIdMap() map[BoardId]string {
	m := make(map[BoardId]string)
	for _, id := range AllBoardIds() {
		m[id] = id.GetNameEn()
	}
	return m
}
