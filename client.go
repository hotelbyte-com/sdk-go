// Package hotelbyte provides a high-quality Go SDK for the HotelByte API
package hotelbyte

import (
	"context"
	"fmt"
	"time"

	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

// Client is the main HotelByte API client
type Client struct {
	config      *Config
	transport   *Transport
	token       string
	tokenExpiry time.Time
}

func (s *Client) Key() string {
	return s.config.Credentials.AppKey
}

// SetToken manually sets the authentication token (for testing purposes)
// This allows setting a token obtained from /api/auth/login for portal users
func (s *Client) SetToken(token string, ttlSeconds int) {
	s.token = token
	s.tokenExpiry = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
}

// Config represents client configuration
type Config struct {
	BaseURL     string
	Credentials Credentials
	HTTPConfig  HTTPConfig
	RetryConfig RetryConfig
}

// Credentials represents authentication credentials
type Credentials struct {
	AppKey    string
	AppSecret string
}

// HTTPConfig represents HTTP configuration
type HTTPConfig struct {
	Timeout         time.Duration
	MaxIdleConns    int
	MaxConnsPerHost int
	UserAgent       string
}

// RetryConfig represents retry configuration
type RetryConfig struct {
	MaxRetries    int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
}

// NewClient creates a new HotelByte client
func NewClient(options ...ClientOption) (*Client, error) {
	config := DefaultConfig()

	for _, option := range options {
		if err := option(config); err != nil {
			return nil, fmt.Errorf("invalid option: %w", err)
		}
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Create transport layer
	transport, err := NewTransport(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport layer: %w", err)
	}

	client := &Client{
		config:    config,
		transport: transport,
	}

	return client, nil
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		BaseURL: "https://api.hotelbyte.com",
		HTTPConfig: HTTPConfig{
			Timeout:         120 * time.Second,
			MaxIdleConns:    100,
			MaxConnsPerHost: 10,
			UserAgent:       "HotelByte-Go-SDK/0.0.1",
		},
		RetryConfig: RetryConfig{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
		},
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("empty base url")
	}

	// Allow empty credentials for external token injection (e.g., portal users)
	// Credentials are optional if token will be set via SetToken()
	return nil
}

// ClientOption represents a client configuration option
type ClientOption func(*Config) error

// WithBaseURL sets the base URL
func WithBaseURL(url string) ClientOption {
	return func(c *Config) error {
		if url == "" {
			return fmt.Errorf("empty base url")
		}
		c.BaseURL = url
		return nil
	}
}

// WithCredentials sets authentication credentials
func WithCredentials(appKey, appSecret string) ClientOption {
	return func(c *Config) error {
		if appKey == "" || appSecret == "" {
			return fmt.Errorf("empty credentials")
		}
		c.Credentials = Credentials{
			AppKey:    appKey,
			AppSecret: appSecret,
		}
		return nil
	}
}

// WithTimeout sets the timeout duration
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Config) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must > 0")
		}
		c.HTTPConfig.Timeout = timeout
		return nil
	}
}

// WithRetryConfig sets the retry configuration
func WithRetryConfig(maxRetries int, initialDelay, maxDelay time.Duration) ClientOption {
	return func(c *Config) error {
		if maxRetries < 0 {
			return fmt.Errorf("invalid max retries")
		}
		c.RetryConfig = RetryConfig{
			MaxRetries:    maxRetries,
			InitialDelay:  initialDelay,
			MaxDelay:      maxDelay,
			BackoffFactor: 2.0,
		}
		return nil
	}
}

// GetConfig returns the client configuration
func (s *Client) GetConfig() *Config {
	return s.config
}

// Close closes the client
func (s *Client) Close() error {
	if s.transport != nil {
		return s.transport.Close()
	}
	return nil
}

func (s *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	return s.transport.Do(ctx, req)
}

type Request = types.HttpRequest
type Response = types.HttpResponse

const (
	// ErrCodeTokenExpired is the error code for token expired (401)
	ErrCodeTokenExpired = 100000401
)

// doWithAuthRetry executes a function with automatic retry on token expiration (401 error)
// If the request fails with ErrCodeTokenExpired, it will re-authenticate and retry once
func doWithAuthRetry[T any](ctx context.Context, client *Client, fn func() (*T, error)) (*T, error) {
	// Ensure user is authenticated before first attempt
	if err := client.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// First attempt
	result, err := fn()
	if err == nil {
		return result, nil
	}

	// Check if error is due to token expired (401)
	if bizErr, ok := types.UnwrapForBizErr(err); ok && bizErr.Code == ErrCodeTokenExpired {
		fmt.Printf("⚠️  Token expired (ErrorCode=%d), re-authenticating...\n", ErrCodeTokenExpired)

		// Re-authenticate
		if refreshErr := client.RefreshToken(ctx); refreshErr != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", refreshErr)
		}

		fmt.Printf("✅ Token refreshed, retrying request...\n")

		// Retry the request with new token
		result, retryErr := fn()
		if retryErr != nil {
			return nil, fmt.Errorf("request failed after token refresh: %w", retryErr)
		}

		return result, nil
	}

	// Return original error if not a token expiration issue
	return nil, err
}
