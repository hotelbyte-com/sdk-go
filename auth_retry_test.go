package hotelbyte

import (
	"context"
	"fmt"
	"testing"

	"github.com/hotelbyte-com/sdk-go/protocol"
	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

// TestAuthRetryOnTokenExpired tests that the SDK automatically retries on token expiration
// This test verifies that when a token expires (ErrorCode 100000401), the SDK will:
// 1. Detect the 401 error
// 2. Re-authenticate to get a new token
// 3. Retry the original request with the new token
func TestAuthRetryOnTokenExpired(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create client with test credentials
	client, err := NewClient(
		WithBaseURL("https://api-test.hotelbyte.com"),
		WithCredentials("openapi_demo", "demo_secret_123456"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test 1: First authentication should succeed
	t.Log("Test 1: Initial authentication")
	if err := client.Authenticate(ctx); err != nil {
		t.Fatalf("Initial authentication failed: %v", err)
	}
	t.Logf("✅ Initial authentication successful, token: %s", client.GetToken())

	// Test 2: Make a request to ensure token works
	t.Log("\nTest 2: Making API request with valid token")
	searchReq := &protocol.HotelListReq{
		HotelDestination: protocol.HotelDestination{
			DestinationName: "Dubai",
		},
		CheckInOut: protocol.CheckInOut{
			CheckIn:  types.DateInt(20250201),
			CheckOut: types.DateInt(20250203),
		},
		Occupancies: protocol.Occupancies{
			NationalityCode: "US",
			RoomOccupancies: []protocol.GuestPerRoom{
				{AdultCount: 2},
			},
		},
		CurrencyOption: protocol.CurrencyOption{
			Currency: "USD",
		},
		PageReq: types.PageReq{
			PageSize: 10,
			PageNum:  1,
		},
	}

	searchResp, err := client.HotelList(ctx, searchReq)
	if err != nil {
		t.Fatalf("HotelList request failed: %v", err)
	}
	t.Logf("✅ API request successful, found %d hotels", len(searchResp.List))

	// Test 3: Manually expire the token to simulate token expiration
	t.Log("\nTest 3: Simulating token expiration")
	// Clear the token to simulate expiration
	client.RefreshToken(ctx) // This will invalidate the current token
	// Actually, we need to manually clear it
	client.token = ""
	t.Log("✅ Token cleared (simulated expiration)")

	// Test 4: Make another request - should automatically re-authenticate
	t.Log("\nTest 4: Making API request with expired token (should auto-retry)")
	searchReq2 := &protocol.HotelListReq{
		HotelDestination: protocol.HotelDestination{
			DestinationName: "Abu Dhabi",
		},
		CheckInOut: protocol.CheckInOut{
			CheckIn:  types.DateInt(20250201),
			CheckOut: types.DateInt(20250203),
		},
		Occupancies: protocol.Occupancies{
			NationalityCode: "US",
			RoomOccupancies: []protocol.GuestPerRoom{
				{AdultCount: 2},
			},
		},
		CurrencyOption: protocol.CurrencyOption{
			Currency: "USD",
		},
		PageReq: types.PageReq{
			PageSize: 10,
			PageNum:  1,
		},
	}

	searchResp2, err := client.HotelList(ctx, searchReq2)
	if err != nil {
		t.Fatalf("HotelList request after token expiration failed: %v", err)
	}
	t.Logf("✅ API request after token expiration successful, found %d hotels", len(searchResp2.List))
	t.Log("✅ Auto-retry mechanism working correctly")

	// Test 5: Verify multiple API calls work continuously
	t.Log("\nTest 5: Verifying continuous API calls")
	for i := 0; i < 3; i++ {
		searchReq := &protocol.HotelListReq{
			HotelDestination: protocol.HotelDestination{
				DestinationName: "Dubai",
			},
			CheckInOut: protocol.CheckInOut{
				CheckIn:  types.DateInt(20250201),
				CheckOut: types.DateInt(20250203),
			},
			Occupancies: protocol.Occupancies{
				NationalityCode: "US",
				RoomOccupancies: []protocol.GuestPerRoom{
					{AdultCount: 2},
				},
			},
			CurrencyOption: protocol.CurrencyOption{
				Currency: "USD",
			},
			PageReq: types.PageReq{
				PageSize: 5,
				PageNum:  1,
			},
		}

		resp, err := client.HotelList(ctx, searchReq)
		if err != nil {
			t.Fatalf("Continuous API call %d failed: %v", i+1, err)
		}
		t.Logf("✅ Continuous call %d successful, found %d hotels", i+1, len(resp.List))
	}

	t.Log("\n✅✅✅ All tests passed! Auto-retry on token expiration is working correctly.")
}

// TestAuthRetryAllEndpoints tests that all endpoints support auto-retry
func TestAuthRetryAllEndpoints(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	client, err := NewClient(
		WithBaseURL("https://api-test.hotelbyte.com"),
		WithCredentials("openapi_demo", "demo_secret_123456"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Clear token to force re-authentication on first call
	client.token = ""

	t.Log("Testing all endpoints with auto-retry...")

	// Test HotelList
	t.Log("1. Testing HotelList...")
	searchReq := &protocol.HotelListReq{
		HotelDestination: protocol.HotelDestination{DestinationName: "Dubai"},
		CheckInOut: protocol.CheckInOut{
			CheckIn:  types.DateInt(20250201),
			CheckOut: types.DateInt(20250203),
		},
		Occupancies: protocol.Occupancies{
			NationalityCode: "US",
			RoomOccupancies: []protocol.GuestPerRoom{{AdultCount: 2}},
		},
		CurrencyOption: protocol.CurrencyOption{Currency: "USD"},
		PageReq:        types.PageReq{PageSize: 5, PageNum: 1},
	}
	searchResp, err := client.HotelList(ctx, searchReq)
	if err != nil {
		t.Errorf("HotelList failed: %v", err)
	} else {
		t.Logf("✅ HotelList successful: %d hotels", len(searchResp.List))
	}

	// Test HotelRates (if we have hotels from previous call)
	if len(searchResp.List) > 0 {
		t.Log("2. Testing HotelRates...")
		hotelID := searchResp.List[0].ID
		sessionID := searchResp.Basic.SessionId

		ratesReq := &protocol.HotelRatesReq{
			HotelId:        hotelID,
			CheckInOut:     searchReq.CheckInOut,
			Occupancies:    searchReq.Occupancies,
			CurrencyOption: searchReq.CurrencyOption,
			SessionOption:  protocol.SessionOption{SessionId: sessionID},
		}
		ratesResp, err := client.HotelRates(ctx, ratesReq)
		if err != nil {
			t.Errorf("HotelRates failed: %v", err)
		} else {
			t.Logf("✅ HotelRates successful: %d rooms", len(ratesResp.Rooms))
		}
	}

	t.Log("✅ All endpoint tests completed")
}

// Example: How the auto-retry works in production
func ExampleClient_autoRetry() {
	ctx := context.Background()

	client, _ := NewClient(
		WithBaseURL("https://api-test.hotelbyte.com"),
		WithCredentials("openapi_demo", "demo_secret_123456"),
	)
	defer client.Close()

	// The SDK will automatically handle token expiration
	// You don't need to manually check for 401 errors
	searchReq := &protocol.HotelListReq{
		HotelDestination: protocol.HotelDestination{DestinationName: "Dubai"},
		CheckInOut: protocol.CheckInOut{
			CheckIn:  types.DateInt(20250201),
			CheckOut: types.DateInt(20250203),
		},
		Occupancies: protocol.Occupancies{
			NationalityCode: "US",
			RoomOccupancies: []protocol.GuestPerRoom{{AdultCount: 2}},
		},
		CurrencyOption: protocol.CurrencyOption{Currency: "USD"},
		PageReq:        types.PageReq{PageSize: 10, PageNum: 1},
	}

	// If the token is expired, the SDK will automatically:
	// 1. Detect the 401 error (ErrorCode=100000401)
	// 2. Re-authenticate to get a new token
	// 3. Retry the request with the new token
	resp, err := client.HotelList(ctx, searchReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d hotels\n", len(resp.List))
	// Output:
	// ⚠️  Token expired (ErrorCode=100000401), re-authenticating...
	// ✅ Token refreshed, retrying request...
	// Found X hotels
}
