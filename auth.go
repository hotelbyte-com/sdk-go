package hotelbyte

import (
	"context"
	"github.com/hotelbyte-com/sdk-go/protocol"
	"github.com/hotelbyte-com/sdk-go/protocol/types"
	"net/http"
)

// Authenticate performs user authentication
func (s *Client) Authenticate(ctx context.Context) error {
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

	r, err := types.NewResponse[protocol.AuthResp](resp.StatusCode, resp.Body)
	if err != nil {
		return err
	}

	// Save token information (expiry may not be provided by backend)
	s.token = r.Ticket
	return nil
}

// GetToken returns the current authentication token
func (s *Client) GetToken() string {
	return s.token
}

// RefreshToken refreshes the authentication token
func (s *Client) RefreshToken(ctx context.Context) error {
	// Clear current token
	s.token = ""
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
