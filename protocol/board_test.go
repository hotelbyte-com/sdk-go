package protocol

import (
	"testing"
)

func TestBoardIdString(t *testing.T) {
	tests := []struct {
		boardId BoardId
		expect  string
	}{
		{BoardIdRoomOnly, "RO"},
		{BoardIdBedBreakfast, "BB"},
		{BoardIdHalfBoard, "HB"},
		{BoardIdFullBoard, "FB"},
		{BoardIdAllInclusive, "AI"},
		{BoardIdBreakfastIncluded, "BI"},
		{BoardIdLunchOnly, "LU"},
		{BoardIdDinnerOnly, "DI"},
		{BoardIdBreakfastDinner, "BD"},
		{BoardIdBreakfastLunch, "BL"},
		{BoardIdLunchDinner, "LD"},
		{BoardIdAllInclusiveTI, "TI"},
	}

	for _, tt := range tests {
		t.Run(string(tt.boardId), func(t *testing.T) {
			if tt.boardId.String() != tt.expect {
				t.Errorf("String() = %v, want %v", tt.boardId.String(), tt.expect)
			}
		})
	}
}

func TestBoardIdGetNameEn(t *testing.T) {
	tests := []struct {
		boardId BoardId
		expect  string
	}{
		{BoardIdRoomOnly, "Room Only"},
		{BoardIdBedBreakfast, "Bed and Breakfast"},
		{BoardIdHalfBoard, "Half Board"},
		{BoardIdFullBoard, "Full Board"},
		{BoardIdAllInclusive, "All Inclusive"},
		{BoardIdBreakfastIncluded, "Breakfast Included"},
	}

	for _, tt := range tests {
		t.Run(string(tt.boardId), func(t *testing.T) {
			if tt.boardId.GetNameEn() != tt.expect {
				t.Errorf("GetNameEn() = %v, want %v", tt.boardId.GetNameEn(), tt.expect)
			}
		})
	}
}

func TestBoardIdGetNameZh(t *testing.T) {
	tests := []struct {
		boardId BoardId
		expect  string
	}{
		{BoardIdRoomOnly, "仅住宿"},
		{BoardIdBedBreakfast, "含早餐"},
		{BoardIdHalfBoard, "半食宿"},
		{BoardIdFullBoard, "全食宿"},
		{BoardIdAllInclusive, "全包"},
		{BoardIdBreakfastIncluded, "仅早餐"},
	}

	for _, tt := range tests {
		t.Run(string(tt.boardId), func(t *testing.T) {
			if tt.boardId.GetNameZh() != tt.expect {
				t.Errorf("GetNameZh() = %v, want %v", tt.boardId.GetNameZh(), tt.expect)
			}
		})
	}
}

func TestBoardIdIsValid(t *testing.T) {
	validBoardIds := []BoardId{
		BoardIdRoomOnly,
		BoardIdBedBreakfast,
		BoardIdHalfBoard,
		BoardIdFullBoard,
		BoardIdAllInclusive,
		BoardIdBreakfastIncluded,
		BoardIdLunchOnly,
		BoardIdDinnerOnly,
		BoardIdBreakfastDinner,
		BoardIdBreakfastLunch,
		BoardIdLunchDinner,
		BoardIdAllInclusiveTI,
		BoardIdBedBreakfast1,
		BoardIdBedBreakfast2,
		BoardIdBedBreakfast3,
	}

	for _, boardId := range validBoardIds {
		t.Run(string(boardId), func(t *testing.T) {
			if !boardId.IsValid() {
				t.Errorf("IsValid() = false, want true for %v", boardId)
			}
		})
	}

	invalidBoardIds := []BoardId{
		BoardId("XX"),
		BoardId("INVALID"),
		BoardId(""),
	}

	for _, boardId := range invalidBoardIds {
		t.Run(string(boardId), func(t *testing.T) {
			if boardId.IsValid() {
				t.Errorf("IsValid() = true, want false for %v", boardId)
			}
		})
	}
}

func TestAllBoardIds(t *testing.T) {
	allIds := AllBoardIds()
	if len(allIds) == 0 {
		t.Error("AllBoardIds() returned empty slice")
	}

	// Check that all returned IDs are valid
	for _, id := range allIds {
		if !id.IsValid() {
			t.Errorf("AllBoardIds() returned invalid ID: %v", id)
		}
	}
}

func TestBoardIdOccupancyVariants(t *testing.T) {
	tests := []struct {
		boardId BoardId
		variant int
	}{
		{BoardIdBedBreakfast1, 1},
		{BoardIdBedBreakfast2, 2},
		{BoardIdBedBreakfast3, 3},
		{BoardIdHalfBoard1, 1},
		{BoardIdHalfBoard2, 2},
		{BoardIdHalfBoard3, 3},
		{BoardIdFullBoard1, 1},
		{BoardIdFullBoard2, 2},
		{BoardIdFullBoard3, 3},
		{BoardIdAllInclusive1, 1},
		{BoardIdAllInclusive2, 2},
		{BoardIdAllInclusive3, 3},
	}

	for _, tt := range tests {
		t.Run(string(tt.boardId), func(t *testing.T) {
			if !tt.boardId.IsValid() {
				t.Errorf("Occupancy variant %v is invalid", tt.boardId)
			}
			nameEn := tt.boardId.GetNameEn()
			if nameEn == "Unknown" {
				t.Errorf("Occupancy variant %v has unknown name", tt.boardId)
			}
		})
	}
}

func TestBoardIdMapKeys(t *testing.T) {
	boardMap := BoardIdMap()

	// Check that the map is not empty
	if len(boardMap) == 0 {
		t.Error("BoardIdMap() returned empty map")
	}

	// Check that all map entries have non-empty values
	for boardId, name := range boardMap {
		if name == "" {
			t.Errorf("BoardIdMap() has empty name for %v", boardId)
		}
		if !boardId.IsValid() {
			t.Errorf("BoardIdMap() contains invalid BoardId: %v", boardId)
		}
	}

	// Check that all valid BoardIds are in the map
	for _, id := range AllBoardIds() {
		if _, exists := boardMap[id]; !exists {
			t.Errorf("BoardIdMap() is missing entry for %v", id)
		}
	}
}
