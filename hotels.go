package hotelbyte

import (
	"context"
	"errors"
	"fmt"
	"github.com/hotelbyte-com/sdk-go/protocol"
	"github.com/hotelbyte-com/sdk-go/protocol/types"
	"github.com/spf13/cast"
	"net/http"
)

func (s *Client) HotelList(ctx context.Context, req *protocol.HotelListReq) (*protocol.HotelListResp, error) {
	// Ensure user is authenticated
	if err := s.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Build request based on real backend structure
	httpReq := &Request{
		Method: http.MethodPost,
		Path:   "/api/search/hotelList",
		Headers: map[string]string{
			"Authorization": s.GetAuthorizationHeader(),
		},
		Body: req, // Use the entire request structure
	}

	// Send request
	resp, err := s.transport.Do(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("hotel search request failed: %w", err)
	}

	// Check response status
	if resp.StatusCode >= 400 {
		return nil, errors.New(cast.ToString(resp.StatusCode))
	}

	return types.NewResponse[protocol.HotelListResp](resp.Body)
}

func (s *Client) HotelRates(ctx context.Context, req *protocol.HotelRatesReq) (*protocol.HotelRatesResp, error) {
	// Ensure user is authenticated
	if err := s.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Build request
	httpReq := &Request{
		Method: http.MethodPost,
		Path:   "/api/search/hotelRates",
		Headers: map[string]string{
			"Authorization": s.GetAuthorizationHeader(),
			"Session-Id":    req.SessionId,
		},
		Body: req, // Use the entire request structure
	}

	// Send request
	resp, err := s.transport.Do(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("get hotel rates request failed: %w", err)
	}

	// Check response status
	if resp.StatusCode >= 400 {
		return nil, errors.New(cast.ToString(resp.StatusCode))
	}

	return types.NewResponse[protocol.HotelRatesResp](resp.Body)
}

func (s *Client) CheckAvail(ctx context.Context, req *protocol.CheckAvailReq) (*protocol.CheckAvailResp, error) {
	// Ensure user is authenticated
	if err := s.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Build request
	httpReq := &Request{
		Method: http.MethodPost,
		Path:   "/api/search/checkAvail",
		Headers: map[string]string{
			"Authorization": s.GetAuthorizationHeader(),
			"Session-Id":    req.SessionId,
		},
		Body: req, // Use the entire request structure
	}

	// Send request
	resp, err := s.transport.Do(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("get hotel rates request failed: %w", err)
	}

	// Check response status
	if resp.StatusCode >= 400 {
		return nil, errors.New(cast.ToString(resp.StatusCode))
	}

	return types.NewResponse[protocol.CheckAvailResp](resp.Body)
}

func (s *Client) Book(ctx context.Context, req *protocol.BookReq) (*protocol.BookResp, error) {
	// Ensure user is authenticated
	if err := s.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Build request
	httpReq := &Request{
		Method: http.MethodPost,
		Path:   "/api/trade/book",
		Headers: map[string]string{
			"Authorization": s.GetAuthorizationHeader(),
			"Session-Id":    req.SessionId,
		},
		Body: req, // Use the entire request structure
	}

	// Send request
	resp, err := s.transport.Do(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("get hotel rates request failed: %w", err)
	}

	// Check response status
	if resp.StatusCode >= 400 {
		return nil, errors.New(cast.ToString(resp.StatusCode))
	}

	return types.NewResponse[protocol.BookResp](resp.Body)
}

func (s *Client) QueryOrders(ctx context.Context, req *protocol.QueryOrdersReq) (*protocol.QueryOrdersResp, error) {
	// Ensure user is authenticated
	if err := s.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Build request
	httpReq := &Request{
		Method: http.MethodPost,
		Path:   "/api/trade/queryOrders",
		Headers: map[string]string{
			"Authorization": s.GetAuthorizationHeader(),
		},
		Body: req, // Use the entire request structure
	}

	// Send request
	resp, err := s.transport.Do(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("get hotel rates request failed: %w", err)
	}

	// Check response status
	if resp.StatusCode >= 400 {
		return nil, errors.New(cast.ToString(resp.StatusCode))
	}

	return types.NewResponse[protocol.QueryOrdersResp](resp.Body)
}

func (s *Client) Cancel(ctx context.Context, req *protocol.CancelReq) (*protocol.CancelResp, error) {
	// Ensure user is authenticated
	if err := s.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Build request
	httpReq := &Request{
		Method: http.MethodPost,
		Path:   "/api/trade/cancel",
		Headers: map[string]string{
			"Authorization": s.GetAuthorizationHeader(),
		},
		Body: req, // Use the entire request structure
	}

	// Send request
	resp, err := s.transport.Do(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("get hotel rates request failed: %w", err)
	}

	// Check response status
	if resp.StatusCode >= 400 {
		return nil, errors.New(cast.ToString(resp.StatusCode))
	}

	return types.NewResponse[protocol.CancelResp](resp.Body)
}
