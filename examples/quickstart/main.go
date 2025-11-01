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
		//hotelbyte.WithBaseURL("http://localhost:8080"),
		hotelbyte.WithBaseURL("https://api-test.hotelbyte.com"),
		hotelbyte.WithCredentials("hotelbyte_api_demo", "hotelbyte_api_demo"),
		hotelbyte.WithTimeout(120*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	testQueries := []string{
		"",
		"disableSimulator=true&onlyAvailableSupplierHotels=true",
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
	fmt.Printf("Session ID: %s\n", searchResp.Basic.SessionId)

	// Example 2: Get hotel rates
	sop := protocol.SessionOption{SessionId: searchResp.Basic.SessionId}
	for _, hotel := range searchResp.List {
		fmt.Printf("Process hotel: %+v (%v)\n", hotel.Name, hotel.ID)
		fmt.Printf("Location: %+v\n", hotel.LatlngCoordinator.Google)
		fmt.Printf("Min price: %.2f %s\n", hotel.MinPrice.Amount, hotel.MinPrice.Currency)
		if handleHotel(ctx, client, hotel, searchReq, sop) {
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
	for _, room := range ratesResp.Rooms {
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
		// 展示取消政策信息
		fmt.Printf("=== Rate Package Details ===\n")
		fmt.Printf("Rate Package ID: %s\n", rate.RatePkgId)
		fmt.Printf("Price: %.2f %s\n", rate.Rate.NetRate.Amount, rate.Rate.NetRate.Currency)
		
		// 展示取消政策
		fmt.Printf("=== Cancellation Policy ===\n")
		fmt.Printf("Refundable Mode: %s\n", rate.RefundableMode)
		
		switch rate.RefundableMode {
		case protocol.RefundableModeFully:
			fmt.Printf("✅ Free cancellation available\n")
		case protocol.RefundableModePartially:
			fmt.Printf("⚠️  Partial refund available\n")
		case protocol.RefundableModeNo:
			fmt.Printf("❌ Non-refundable\n")
		}
		
		if !rate.RefundableUntil.IsZero() {
			fmt.Printf("Free cancellation until: %s\n", rate.RefundableUntil.Format("2006-01-02 15:04:05 MST"))
		}
		
		if len(rate.CancelFees) > 0 {
			fmt.Printf("Cancellation fees:\n")
			for i, fee := range rate.CancelFees {
				fmt.Printf("  %d. Until %s: %.2f %s\n", 
					i+1, 
					fee.Until.Format("2006-01-02 15:04:05 MST"), 
					fee.Fee.Amount, 
					fee.Fee.Currency)
			}
		}
		
		// 根据取消政策决定是否继续预订
		if rate.RefundableMode == protocol.RefundableModeNo {
			fmt.Printf("⚠️  Skipping non-refundable rate for demo purposes\n")
			continue
		}
		
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
		log.Printf("Check availability failed: %v\n", err)
		return false
	}
	log.Printf("Check Availibility, status:%v\n", checkAvailResp.Status)
	if checkAvailResp.Status != protocol.CheckAvailStatusAvailable {
		return false
	}
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
	fmt.Printf("Booking created successfully! %s\n", Pretty(bookingResp))
	//fmt.Printf("Supplier Order ID: %v\n", bookingResp.HotelOrder.SupplierReferenceNo)
	//fmt.Printf("Order Status: %d\n", bookingResp.HotelOrder.Status)
	//fmt.Printf("Customer Order ID: %v\n", bookingResp.HotelOrder.CustomerReferenceNo)

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

	// Example 5: Cancel booking with detailed cancellation policy information
	fmt.Println("=== Cancelling booking ===")
	fmt.Printf("Original cancellation policy for this rate:\n")
	fmt.Printf("- Refundable Mode: %s\n", rate.RefundableMode)
	
	if len(rate.CancelFees) > 0 {
		fmt.Printf("- Expected cancellation fees based on booking time:\n")
		now := time.Now()
		for i, fee := range rate.CancelFees {
			if now.Before(fee.Until) {
				fmt.Printf("  Current fee (until %s): %.2f %s\n", 
					fee.Until.Format("2006-01-02 15:04"), 
					fee.Fee.Amount, 
					fee.Fee.Currency)
				break
			} else if i == len(rate.CancelFees)-1 {
				fmt.Printf("  Current fee: %.2f %s (final penalty)\n", 
					fee.Fee.Amount, 
					fee.Fee.Currency)
			}
		}
	} else {
		fmt.Printf("- No cancellation fees expected\n")
	}
	
	cancelResp, err := client.Cancel(ctx, &protocol.CancelReq{
		CustomerReferenceNo: bookingReq.CustomerReferenceNo,
		TestOption:          top,
	})
	if err != nil {
		log.Printf("❌ Cancel booking failed: %v", err)
		return false
	}
	
	fmt.Printf("✅ Cancel booking successfully!\n")
	fmt.Printf("- Final status: %+v\n", cancelResp.Status)
	fmt.Printf("- Actual service fee charged: %.2f %s\n", 
		cancelResp.ServiceFee.Amount, 
		cancelResp.ServiceFee.Currency)
	
	// 比较预期费用和实际费用
	if len(rate.CancelFees) > 0 {
		now := time.Now()
		for _, fee := range rate.CancelFees {
			if now.Before(fee.Until) {
				if fee.Fee.Amount != cancelResp.ServiceFee.Amount {
					fmt.Printf("⚠️  Note: Expected fee (%.2f) differs from actual fee (%.2f)\n", 
						fee.Fee.Amount, cancelResp.ServiceFee.Amount)
				} else {
					fmt.Printf("✅ Cancellation fee matches expectation\n")
				}
				break
			}
		}
	}
	return true
}

func Pretty(v interface{}) string {
	s, _ := sonic.MarshalIndent(v, "", "  ")
	return string(s)
}
