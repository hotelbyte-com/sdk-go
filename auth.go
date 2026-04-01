package hotelbyte

import (
	"context"
	"fmt"
	"time"

	"net/http"

	"github.com/hotelbyte-com/sdk-go/protocol"
	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

const defaultAuthTTLSeconds int64 = 24 * 3600

// Authenticate performs user authentication
func (s *Client) Authenticate(ctx context.Context) error {
	// 如果 token 存在且未过期（提前 5 分钟刷新），直接返回
	s.mu.RLock()
	valid := s.token != "" && time.Now().Before(s.tokenExpiry.Add(-5*time.Minute))
	s.mu.RUnlock()
	if valid {
		return nil
	}

	_, err, _ := s.sf.Do("auth", func() (interface{}, error) {
		// double check inside singleflight
		s.mu.RLock()
		valid = s.token != "" && time.Now().Before(s.tokenExpiry.Add(-5*time.Minute))
		s.mu.RUnlock()
		if valid {
			return nil, nil
		}

		token, expiry, err := s.fetchToken(ctx, defaultAuthTTLSeconds)
		if err != nil {
			return nil, err
		}
		
		s.mu.Lock()
		s.token = token
		s.tokenExpiry = expiry
		s.mu.Unlock()
		return nil, nil
	})
	return err
}

func (s *Client) fetchToken(ctx context.Context, ttlSeconds int64) (string, time.Time, error) {
	if ttlSeconds <= 0 {
		ttlSeconds = defaultAuthTTLSeconds
	}

	// Build authentication request
	req := &protocol.AuthReq{
		AppKey:    s.config.Credentials.AppKey,
		AppSecret: s.config.Credentials.AppSecret,
		TTL:       ttlSeconds,
	}

	httpReq := &Request{
		Method: http.MethodPost,
		Path:   "/api/auth/ticket",
		Body:   req,
	}
	resp, err := s.transport.Do(ctx, httpReq)
	if err != nil {
		// On transport error, try fallback path if available
		return "", time.Time{}, err
	}

	r, err := types.NewResponseData[protocol.AuthResp](resp)
	if err != nil {
		return "", time.Time{}, err
	}
	if r.Ticket == "" {
		return "", time.Time{}, fmt.Errorf("empty ticket from auth response")
	}

	expiry := time.Now().Add(time.Duration(req.TTL) * time.Second)
	return r.Ticket, expiry, nil
}

// GetToken returns the current authentication token
func (s *Client) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token
}

// SetToken sets an external authentication token (for example, JWT from portal login)
// This allows using a token obtained from /api/auth/login without re-authentication
// The token will be used directly in Authorization headers for subsequent requests
// If ttlSeconds is 0 or negative, it will set a far future expiry (365 days) for external tokens
func (s *Client) SetToken(token string, ttlSeconds ...int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
	if len(ttlSeconds) > 0 && ttlSeconds[0] > 0 {
		s.tokenExpiry = time.Now().Add(time.Duration(ttlSeconds[0]) * time.Second)
	} else {
		// Set a far future expiry for external tokens (typically JWT with own expiry)
		// The API will reject expired tokens anyway
		s.tokenExpiry = time.Now().Add(365 * 24 * time.Hour)
	}
}

// RefreshToken refreshes the authentication token
func (s *Client) RefreshToken(ctx context.Context) error {
	_, err, _ := s.sf.Do("refresh", func() (interface{}, error) {
		// Keep old token until refresh succeeds, avoiding empty-token window under concurrency.
		token, expiry, err := s.fetchToken(ctx, defaultAuthTTLSeconds)
		if err != nil {
			return nil, err
		}
		s.mu.Lock()
		s.token = token
		s.tokenExpiry = expiry
		s.mu.Unlock()
		return nil, nil
	})
	return err
}

// GetAuthToken returns the current token (alias for GetToken for backward compatibility)
func (s *Client) GetAuthToken() string {
	return s.GetToken()
}

// GetAuthorizationHeader returns the authorization header value
func (s *Client) GetAuthorizationHeader() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.token == "" {
		return ""
	}
	return "Bearer " + s.token
}
