package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"

	"github.com/spf13/cast"

	hotelbyte "github.com/hotelbyte-com/sdk-go"
	"github.com/hotelbyte-com/sdk-go/protocol"
	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

func main() {
	// Initialize SDK client with credentials (use client options API)
	client, err := hotelbyte.NewClient(
		hotelbyte.WithBaseURL("http://localhost:8080"),
		hotelbyte.WithCredentials("hotelbyte_api_demo", "hotelbyte_api_demo"),
		hotelbyte.WithTimeout(120*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	testQueries := []string{
		//"",
		"disableSimulator=true&onlyAvailableHotels=true",
	}
	for _, tq := range testQueries {
		run(client, tq)
	}
}

func run(client *hotelbyte.Client, tq string) {
	ctx := context.Background()

	// Example 1: Search for hotels
	fmt.Println("=== Searching for hotels ===")
	searchReq := &protocol.HotelListReq{
		HotelDestination: protocol.HotelDestination{
			DestinationName: "Dubai",
		},
		CheckInOut: protocol.CheckInOut{
			CheckIn:  20250315, // YYYYMMDD format
			CheckOut: 20250317, // YYYYMMDD format
		},
		Occupancies: protocol.Occupancies{
			NationalityCode: "US",
			RoomOccupancies: []protocol.GuestPerRoom{
				{
					AdultCount: 2,
				},
				{
					AdultCount: 1,
				},
			},
		},
		CurrencyOption: protocol.CurrencyOption{
			Currency: "USD",
		},
		MaxRatesPerHotel: 3,
		SortBy:           "price-asc",
		PageReq: types.PageReq{
			PageSize: 1000000,
			PageNum:  1,
		},
		TestOption: protocol.TestOption{
			Test: tq,
		},
	}

	searchResp, err := client.HotelList(ctx, searchReq)
	if err != nil {
		log.Fatalf("Hotel search failed: %v", err)
	}

	fmt.Printf("Found %d hotels\n", len(searchResp.List))
	if len(searchResp.List) > 0 {
		hotel := searchResp.List[0]
		fmt.Printf("First hotel: %+v\n", hotel.Name)
		fmt.Printf("Location: %+v\n", hotel.LatlngCoordinator.Google)
		fmt.Printf("Min price: %.2f %s\n", hotel.MinPrice.Amount, hotel.MinPrice.Currency)
	}
	fmt.Printf("Session ID: %s\n", searchResp.Basic.SessionId)

	// Example 2: Get hotel rates
	sop := protocol.SessionOption{SessionId: searchResp.Basic.SessionId}
	for _, h := range searchResp.List {
		if handleHotel(ctx, client, h, searchReq, sop) {
			break
		}
	}
}

func handleHotel(ctx context.Context, client *hotelbyte.Client, hotel *protocol.Hotel, searchReq *protocol.HotelListReq, sop protocol.SessionOption) bool {
	fmt.Println("=== Getting hotel rates ===")
	ratesReq := &protocol.HotelRatesReq{
		HotelId: hotel.ID,
		CheckInOut: protocol.CheckInOut{
			CheckIn:  types.NewDateIntFromTime(time.Now().AddDate(0, 0, 1)),
			CheckOut: types.NewDateIntFromTime(time.Now().AddDate(0, 0, 3)),
		},
		Occupancies:    searchReq.Occupancies,
		CurrencyOption: searchReq.CurrencyOption,
		SessionOption:  sop,
		TestOption:     searchReq.TestOption,
	}

	ratesResp, err := client.HotelRates(ctx, ratesReq)
	if err != nil {
		log.Fatalf("Get rates failed: %v", err)
	}

	fmt.Printf("Found %d rooms with rates\n", len(ratesResp.Rooms))
	if len(ratesResp.Rooms) > 0 {
		room := ratesResp.Rooms[0]
		fmt.Printf("Room: %+v\n", room)
		fmt.Printf("Available rates: %d\n", len(room.Rates))
		if handleRoom(ctx, client, room, sop, searchReq.TestOption) {
			return true
		}
	}
	return false
}

func handleRoom(ctx context.Context, client *hotelbyte.Client, room *protocol.Room, sop protocol.SessionOption, top protocol.TestOption) bool {
	for _, rate := range room.Rates {
		if handleRate(ctx, client, rate, sop, top) {
			return true
		}
	}
	return false
}

func handleRate(ctx context.Context, client *hotelbyte.Client, rate protocol.RoomRatePkg, sop protocol.SessionOption, top protocol.TestOption) bool {
	fmt.Println("=== Check Availability ===")
	checkAvailReq := &protocol.CheckAvailReq{
		RatePkgId:     rate.RatePkgId,
		SessionOption: sop,
	}
	checkAvailResp, err := client.CheckAvail(ctx, checkAvailReq)
	if err != nil {
		log.Fatalf("Check availability failed: %v", err)
	}
	log.Printf("Check Availibility, status:%v\n", checkAvailResp.Status)

	fmt.Println("=== Creating a booking ===")
	bookingReq := &protocol.BookReq{
		CustomerReferenceNo: cast.ToString(time.Now().Unix()),
		RatePkgId:           checkAvailReq.RatePkgId,
		TestOption:          top,
		Holder: protocol.Holder{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone: protocol.Phone{
				CountryCode:   "US",
				CountryNumber: 1,
				Number:        "5551234567",
			},
		},
		Guests: []protocol.Guest{
			{
				FirstName:       "John",
				LastName:        "Doe",
				RoomIndex:       1,
				NationalityCode: "US",
				Age:             35,
			},
			{
				FirstName:       "Jane",
				LastName:        "Doe",
				RoomIndex:       1,
				NationalityCode: "US",
				Age:             32,
			},
		},
		SessionOption: sop,
	}

	bookingResp, err := client.Book(ctx, bookingReq)
	if err != nil {
		log.Printf("Booking creation failed: %v", err)
		panic(err)
	}
	fmt.Printf("Booking created successfully!\n")
	fmt.Printf("Supplier Order ID: %v\n", bookingResp.HotelOrder.SupplierReferenceNo)
	fmt.Printf("Order Status: %d\n", bookingResp.HotelOrder.Status)
	fmt.Printf("Customer Order ID: %v\n", bookingResp.HotelOrder.CustomerReferenceNo)

	// Example 4: Query booking details
	fmt.Println("=== Getting booking details ===")
	getBookingReq := &protocol.QueryOrdersReq{
		CustomerReferenceNos: []string{bookingReq.CustomerReferenceNo},
		TestOption:           top,
	}

	bookingDetails, err := client.QueryOrders(ctx, getBookingReq)
	if err != nil {
		log.Printf("Get booking details failed: %v", err)
	} else {
		log.Printf("Booking Details: %s\n", Pretty(bookingDetails.Orders))
	}

	cancelResp, err := client.Cancel(ctx, &protocol.CancelReq{
		CustomerReferenceNo: bookingReq.CustomerReferenceNo,
		TestOption:          top,
	})
	if err != nil {
		log.Printf("Cancel booking failed: %v", err)
		panic(err)
	}
	log.Printf("Cancel booking successfully, status:%+v, serviceFee:%+v\n", cancelResp.Status, cancelResp.ServiceFee)
	return true
}

func Pretty(v interface{}) string {
	s, _ := sonic.MarshalIndent(v, "", "  ")
	return string(s)
}
