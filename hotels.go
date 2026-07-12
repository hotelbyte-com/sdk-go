package hotelbyte

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hotelbyte-com/sdk-go/protocol"
	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

func (s *Client) HotelList(ctx context.Context, req *protocol.HotelListReq) (*protocol.HotelListResp, error) {
	return doWithAuthRetry(ctx, s, func() (*protocol.HotelListResp, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodPost,
			Path:   "/api/search/hotelList",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
			Body: req,
		}
		if req.Test != "" {
			httpReq.Headers["Test"] = req.Test
		}
		if req.Currency != "" {
			httpReq.Headers["Currency"] = req.Currency
		}

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("hotel search request failed: %w", err)
		}

		return types.NewResponseData[protocol.HotelListResp](resp)
	})
}

func (s *Client) HotelRates(ctx context.Context, req *protocol.HotelRatesReq) (*protocol.HotelRatesResp, error) {
	return doWithAuthRetry(ctx, s, func() (*protocol.HotelRatesResp, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodPost,
			Path:   "/api/search/hotelRates",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
			Body: req,
		}
		setOptionalHeader(httpReq.Headers, "Session-Id", req.SessionId)
		setOptionalHeader(httpReq.Headers, "Test", req.Test)
		setOptionalHeader(httpReq.Headers, "Currency", req.Currency)

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("get hotel rates request failed: %w", err)
		}

		return types.NewResponseData[protocol.HotelRatesResp](resp)
	})
}

func (s *Client) CheckAvail(ctx context.Context, req *protocol.CheckAvailReq) (*protocol.CheckAvailResp, error) {
	return doWithAuthRetry(ctx, s, func() (*protocol.CheckAvailResp, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodPost,
			Path:   "/api/search/checkAvail",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
			Body: req,
		}
		setOptionalHeader(httpReq.Headers, "Session-Id", req.SessionId)
		setOptionalHeader(httpReq.Headers, "Test", req.Test)

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("check avail request failed: %w", err)
		}

		return types.NewResponseData[protocol.CheckAvailResp](resp)
	})
}

func (s *Client) Book(ctx context.Context, req *protocol.BookReq) (*protocol.BookResp, error) {
	return doWithAuthRetry(ctx, s, func() (*protocol.BookResp, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodPost,
			Path:   "/api/trade/book",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
			Body: req,
		}
		setOptionalHeader(httpReq.Headers, "Session-Id", req.SessionId)
		setOptionalHeader(httpReq.Headers, "Test", req.Test)

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("book request failed: %w", err)
		}

		return types.NewResponseData[protocol.BookResp](resp)
	})
}

func (s *Client) QueryOrders(ctx context.Context, req *protocol.QueryOrdersReq) (*protocol.QueryOrdersResp, error) {
	return doWithAuthRetry(ctx, s, func() (*protocol.QueryOrdersResp, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodPost,
			Path:   "/api/trade/queryOrders",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
			Body: req,
		}
		setOptionalHeader(httpReq.Headers, "Test", req.Test)

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("query orders request failed: %w", err)
		}

		return types.NewResponseData[protocol.QueryOrdersResp](resp)
	})
}

func (s *Client) Cancel(ctx context.Context, req *protocol.CancelReq) (*protocol.CancelResp, error) {
	return doWithAuthRetry(ctx, s, func() (*protocol.CancelResp, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodPost,
			Path:   "/api/trade/cancel",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
			Body: req,
		}
		setOptionalHeader(httpReq.Headers, "Test", req.Test)

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("cancel request failed: %w", err)
		}

		return types.NewResponseData[protocol.CancelResp](resp)
	})
}

// setOptionalHeader sets a header only when the value is non-empty.
// Sending empty-string headers (e.g. "Test: ") can cause the backend to
// enter a different code path (e.g. static-only mode without sessionId).
func setOptionalHeader(headers map[string]string, key, value string) {
	if value != "" {
		headers[key] = value
	}
}
