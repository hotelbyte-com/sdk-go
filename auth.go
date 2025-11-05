package hotelbyte

import (
	"context"
	"time"

	"net/http"

	"github.com/hotelbyte-com/sdk-go/protocol"
	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

// Authenticate performs user authentication
func (s *Client) Authenticate(ctx context.Context) error {
	// 如果 token 存在且未过期（提前 5 分钟刷新），直接返回
	if s.token != "" && time.Now().Before(s.tokenExpiry.Add(-5*time.Minute)) {
		return nil
	}

	// Build authentication request
	req := &protocol.AuthReq{
		AppKey:    s.config.Credentials.AppKey,
		AppSecret: s.config.Credentials.AppSecret,
		TTL:       24 * 3600,
	}

	httpReq := &Request{
		Method: http.MethodPost,
		Path:   "/api/auth/ticket",
		Body:   req,
	}
	resp, err := s.transport.Do(ctx, httpReq)
	if err != nil {
		// On transport error, try fallback path if available
		return err
	}

	r, err := types.NewResponseData[protocol.AuthResp](resp)
	if err != nil {
		return err
	}

	// 保存 token 和过期时间
	s.token = r.Ticket
	s.tokenExpiry = time.Now().Add(time.Duration(req.TTL) * time.Second)
	return nil
}

// GetToken returns the current authentication token
func (s *Client) GetToken() string {
	return s.token
}

// RefreshToken refreshes the authentication token
func (s *Client) RefreshToken(ctx context.Context) error {
	// Clear current token and expiry to force re-authentication
	s.token = ""
	s.tokenExpiry = time.Time{}
	// Re-authenticate
	return s.Authenticate(ctx)
}

// GetAuthToken returns the current token (alias for GetToken for backward compatibility)
func (s *Client) GetAuthToken() string {
	return s.GetToken()
}

// GetAuthorizationHeader returns the authorization header value
func (s *Client) GetAuthorizationHeader() string {
	if s.token == "" {
		return ""
	}
	return "Bearer " + s.token
}
