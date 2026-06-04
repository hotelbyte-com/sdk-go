package hotelbyte

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestRoomMappingLiveSmoke(t *testing.T) {
	baseURL := os.Getenv("HOTELBYTE_ROOM_MAPPING_LIVE_BASE_URL")
	appKey := os.Getenv("HOTELBYTE_ROOM_MAPPING_LIVE_APP_KEY")
	appSecret := os.Getenv("HOTELBYTE_ROOM_MAPPING_LIVE_APP_SECRET")
	if baseURL == "" || appKey == "" || appSecret == "" {
		t.Skip("set HOTELBYTE_ROOM_MAPPING_LIVE_BASE_URL, HOTELBYTE_ROOM_MAPPING_LIVE_APP_KEY, and HOTELBYTE_ROOM_MAPPING_LIVE_APP_SECRET")
	}

	client, err := NewClient(
		WithBaseURL(baseURL),
		WithCredentials(appKey, appSecret),
		WithTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mapping, err := client.MapRooms(ctx, &RoomMappingRequest{
		CountryCode: "PT",
		List: []RoomMappingRoomInfo{
			{HotelID: "H1", Supplier: "HB", RoomCode1: "A", Amount: "100", Name: "Deluxe King Room"},
			{HotelID: "H1", Supplier: "HB", RoomCode1: "B", Amount: "102", Name: "Deluxe King Room"},
		},
	})
	if err != nil {
		t.Fatalf("map rooms: %v", err)
	}
	if mapping.TenantID == "" || mapping.MappingID == "" || mapping.Stats.TotalRooms != 2 {
		t.Fatalf("unexpected mapping response: %+v", mapping)
	}

	report, err := client.RoomMappingUsageReport(ctx)
	if err != nil {
		t.Fatalf("usage report: %v", err)
	}
	if report.TenantID != mapping.TenantID || report.Requests < 1 || report.RoomUnits < 2 {
		t.Fatalf("unexpected usage report: %+v", report)
	}
}
