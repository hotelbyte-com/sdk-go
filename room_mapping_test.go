package hotelbyte

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoomMappingV1LookupAndCorrection(t *testing.T) {
	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/mappings/rm_1":
			_ = json.NewEncoder(w).Encode(RoomMappingResponse{
				MappingID: "rm_1",
				Stats:     RoomMappingStats{TotalRooms: 2},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/rooms/correct":
			_ = json.NewEncoder(w).Encode(RoomMappingCorrectionResponse{
				MappingID: "rm_1",
				Applied:   1,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/usage/report":
			_ = json.NewEncoder(w).Encode(RoomMappingUsageReport{
				TenantID:             "tenant-a",
				Requests:             1,
				RoomUnits:            2,
				FreeRoomUnits:        1,
				BillableRoomUnits:    1,
				RoomMappingUnitPrice: 0.25,
				Currency:             "USD",
				EstimatedAmount:      0.25,
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithDisableAutoAuthRetry(true),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	mapping, err := client.GetRoomMapping(context.Background(), "rm_1")
	if err != nil {
		t.Fatalf("get mapping: %v", err)
	}
	if mapping.MappingID != "rm_1" || mapping.Stats.TotalRooms != 2 {
		t.Fatalf("unexpected mapping: %+v", mapping)
	}

	correction, err := client.CorrectRooms(context.Background(), &RoomMappingCorrectionRequest{
		MappingID: "rm_1",
		Corrections: []RoomMappingCorrection{
			{RoomCode1: "B", Group: "Manual"},
		},
	})
	if err != nil {
		t.Fatalf("correct rooms: %v", err)
	}
	if correction.Applied != 1 {
		t.Fatalf("applied = %d, want 1", correction.Applied)
	}

	if len(seen) != 2 || seen[0] != "GET /v1/mappings/rm_1" || seen[1] != "POST /v1/rooms/correct" {
		t.Fatalf("unexpected request paths: %#v", seen)
	}
}

func TestMapRoomsUsesV1Endpoint(t *testing.T) {
	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost && r.URL.Path == "/v1/rooms/map" {
			_ = json.NewEncoder(w).Encode(RoomMappingResponse{
				MappingID: "rm_1",
				Stats:     RoomMappingStats{TotalRooms: 2},
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithDisableAutoAuthRetry(true),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	mapping, err := client.MapRooms(context.Background(), &RoomMappingRequest{
		CountryCode: "PT",
		List: []RoomMappingRoomInfo{
			{HotelID: "H1", Supplier: "HB", RoomCode1: "A", Name: "Deluxe King"},
			{HotelID: "H1", Supplier: "HB", RoomCode1: "B", Name: "Deluxe King"},
		},
	})
	if err != nil {
		t.Fatalf("map rooms: %v", err)
	}
	if mapping.MappingID != "rm_1" || mapping.Stats.TotalRooms != 2 {
		t.Fatalf("unexpected mapping: %+v", mapping)
	}
	if len(seen) != 1 || seen[0] != "POST /v1/rooms/map" {
		t.Fatalf("unexpected request paths: %#v", seen)
	}
}

func TestRoomMappingUsageReport(t *testing.T) {
	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/v1/usage/report" {
			_ = json.NewEncoder(w).Encode(RoomMappingUsageReport{
				TenantID:             "tenant-a",
				Requests:             1,
				RoomUnits:            2,
				FreeRoomUnits:        1,
				BillableRoomUnits:    1,
				RoomMappingUnitPrice: 0.25,
				Currency:             "USD",
				EstimatedAmount:      0.25,
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithDisableAutoAuthRetry(true),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	report, err := client.RoomMappingUsageReport(context.Background())
	if err != nil {
		t.Fatalf("usage report: %v", err)
	}
	if report.TenantID != "tenant-a" || report.BillableRoomUnits != 1 || report.EstimatedAmount != 0.25 {
		t.Fatalf("unexpected usage report: %+v", report)
	}
	if len(seen) != 1 || seen[0] != "GET /v1/usage/report" {
		t.Fatalf("unexpected request paths: %#v", seen)
	}
}
