package hotelbyte

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

type RoomMappingRoomInfo struct {
	HotelID   string `json:"hotelID"`
	Supplier  string `json:"supplier"`
	RoomCode1 string `json:"roomCode1,omitempty"`
	RoomCode2 string `json:"roomCode2,omitempty"`
	Amount    string `json:"amount"`
	Name      string `json:"name"`
}

type RoomMappingRequest struct {
	CountryCode string                `json:"countryCode"`
	List        []RoomMappingRoomInfo `json:"list"`
}

type RoomMappingMappedRoomInfo struct {
	HotelID    string  `json:"hotelID"`
	Supplier   string  `json:"supplier"`
	RoomCode1  string  `json:"roomCode1"`
	RoomCode2  string  `json:"roomCode2"`
	Amount     float64 `json:"amount"`
	Name       string  `json:"name"`
	Group      string  `json:"group"`
	Confidence float64 `json:"confidence"`
	Score      float64 `json:"score"`
}

type RoomMappingStats struct {
	TotalRooms    int     `json:"totalRooms"`
	GroupCount    int     `json:"groupCount"`
	AvgGroupSize  float64 `json:"avgGroupSize"`
	UnmappedCount int     `json:"unmappedCount"`
}

type RoomMappingResponse struct {
	MappingID string                                 `json:"mappingId,omitempty"`
	TenantID  string                                 `json:"tenantId,omitempty"`
	Groups    map[string][]RoomMappingMappedRoomInfo `json:"groups"`
	Stats     RoomMappingStats                       `json:"stats"`
}

type RoomMappingBatchRequest struct {
	Requests []RoomMappingRequest `json:"requests"`
}

type RoomMappingBatchResponse struct {
	Results []RoomMappingResponse `json:"results"`
}

type RoomMappingCorrection struct {
	HotelID    string   `json:"hotelID,omitempty"`
	Supplier   string   `json:"supplier,omitempty"`
	RoomCode1  string   `json:"roomCode1,omitempty"`
	RoomCode2  string   `json:"roomCode2,omitempty"`
	Group      string   `json:"group"`
	Confidence *float64 `json:"confidence,omitempty"`
}

type RoomMappingCorrectionRequest struct {
	MappingID   string                  `json:"mappingId"`
	Corrections []RoomMappingCorrection `json:"corrections"`
}

type RoomMappingCorrectionResponse struct {
	TenantID  string              `json:"tenantId,omitempty"`
	MappingID string              `json:"mappingId"`
	Result    RoomMappingResponse `json:"result"`
	Quality   map[string]any      `json:"quality,omitempty"`
	Applied   int                 `json:"applied"`
}

type RoomMappingUsageReport struct {
	TenantID             string  `json:"tenantId"`
	Requests             int64   `json:"requests"`
	RoomUnits            int64   `json:"roomUnits"`
	FreeRoomUnits        int64   `json:"freeRoomUnits"`
	BillableRoomUnits    int64   `json:"billableRoomUnits"`
	RoomMappingUnitPrice float64 `json:"roomMappingUnitPrice"`
	Currency             string  `json:"currency"`
	EstimatedAmount      float64 `json:"estimatedAmount"`
	Scope                string  `json:"scope,omitempty"`
}

func (s *Client) MapRooms(ctx context.Context, req *RoomMappingRequest) (*RoomMappingResponse, error) {
	return doWithAuthRetry(ctx, s, func() (*RoomMappingResponse, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodPost,
			Path:   "/v1/rooms/map",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
			Body: req,
		}

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("map rooms request failed: %w", err)
		}
		return decodeWrappedOrRaw[RoomMappingResponse](resp)
	})
}

func (s *Client) MapRoomsBatch(ctx context.Context, req *RoomMappingBatchRequest) (*RoomMappingBatchResponse, error) {
	return doWithAuthRetry(ctx, s, func() (*RoomMappingBatchResponse, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodPost,
			Path:   "/api/mapping/rooms/batch",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
			Body: req,
		}

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("map rooms batch request failed: %w", err)
		}
		return decodeWrappedOrRaw[RoomMappingBatchResponse](resp)
	})
}

func (s *Client) GetRoomMapping(ctx context.Context, mappingID string) (*RoomMappingResponse, error) {
	return doWithAuthRetry(ctx, s, func() (*RoomMappingResponse, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodGet,
			Path:   "/v1/mappings/" + url.PathEscape(mappingID),
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
		}

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("get room mapping request failed: %w", err)
		}
		return decodeWrappedOrRaw[RoomMappingResponse](resp)
	})
}

func (s *Client) CorrectRooms(ctx context.Context, req *RoomMappingCorrectionRequest) (*RoomMappingCorrectionResponse, error) {
	return doWithAuthRetry(ctx, s, func() (*RoomMappingCorrectionResponse, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodPost,
			Path:   "/v1/rooms/correct",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
			Body: req,
		}

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("correct rooms request failed: %w", err)
		}
		return decodeWrappedOrRaw[RoomMappingCorrectionResponse](resp)
	})
}

func (s *Client) RoomMappingUsageReport(ctx context.Context) (*RoomMappingUsageReport, error) {
	return doWithAuthRetry(ctx, s, func() (*RoomMappingUsageReport, error) {
		httpReq := &types.HttpRequest{
			Method: http.MethodGet,
			Path:   "/v1/usage/report",
			Headers: map[string]string{
				"Authorization": s.GetAuthorizationHeader(),
			},
		}

		resp, err := s.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, fmt.Errorf("room mapping usage report request failed: %w", err)
		}
		return decodeWrappedOrRaw[RoomMappingUsageReport](resp)
	})
}

func decodeWrappedOrRaw[T any](resp *types.HttpResponse) (*T, error) {
	if resp == nil || len(resp.Body) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	wrapped, wrappedErr := types.NewResponseData[T](resp)
	if wrappedErr == nil && wrapped != nil {
		return wrapped, nil
	}

	var raw T
	if err := json.Unmarshal(resp.Body, &raw); err != nil {
		return nil, wrappedErr
	}
	return &raw, nil
}
