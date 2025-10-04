package hotelbyte

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Transport represents HTTP transport layer
type Transport struct {
	client *http.Client
	config *Config
}

// NewTransport creates a new transport layer
func NewTransport(config *Config) (*Transport, error) {
	client := &http.Client{
		Timeout: config.HTTPConfig.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        config.HTTPConfig.MaxIdleConns,
			MaxIdleConnsPerHost: config.HTTPConfig.MaxConnsPerHost,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &Transport{
		client: client,
		config: config,
	}, nil
}

// Request represents HTTP request
type Request struct {
	Method  string
	Path    string
	Query   url.Values
	Headers map[string]string
	Body    interface{}
}

// Response represents HTTP response
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// Do executes HTTP request
func (t *Transport) Do(ctx context.Context, req *Request) (*Response, error) {
	// Build complete URL
	fullURL, err := t.buildURL(req.Path, req.Query)
	if err != nil {
		return nil, fmt.Errorf("构建 URL 失败: %w", err)
	}

	// Serialize request body
	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := sonic.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create HTTP request failed: %w", err)
	}

	// Set default headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", t.config.HTTPConfig.UserAgent)

	// Set custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request with retry mechanism
	v, err := t.doWithRetry(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	logrus.Infof("response headers : %v", v.Headers) // for debug
	return v, nil
}

// buildURL builds complete URL
func (t *Transport) buildURL(path string, query url.Values) (string, error) {
	baseURL, err := url.Parse(t.config.BaseURL)
	if err != nil {
		return "", err
	}

	pathURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	fullURL := baseURL.ResolveReference(pathURL)
	if query != nil {
		fullURL.RawQuery = query.Encode()
	}

	return fullURL.String(), nil
}

// doWithRetry executes request with retry mechanism
func (t *Transport) doWithRetry(ctx context.Context, req *http.Request) (*Response, error) {
	var lastErr error

	for attempt := 0; attempt <= t.config.RetryConfig.MaxRetries; attempt++ {
		// If not the first attempt, wait for a while
		if attempt > 0 {
			delay := t.calculateDelay(attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		// Execute request
		resp, err := t.client.Do(req)
		if err != nil {
			lastErr = err
			if !t.shouldRetry(err) {
				break
			}
			continue
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		response := &Response{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       body,
		}

		// Check if retry is needed
		if t.shouldRetryResponse(response) && attempt < t.config.RetryConfig.MaxRetries {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
			continue
		}

		return response, nil
	}

	return nil, fmt.Errorf("request failed, retried %d times: %w", t.config.RetryConfig.MaxRetries, lastErr)
}

// calculateDelay calculates retry delay
func (t *Transport) calculateDelay(attempt int) time.Duration {
	delay := time.Duration(float64(t.config.RetryConfig.InitialDelay) *
		pow(t.config.RetryConfig.BackoffFactor, float64(attempt-1)))

	if delay > t.config.RetryConfig.MaxDelay {
		delay = t.config.RetryConfig.MaxDelay
	}

	return delay
}

// shouldRetry determines if network error should be retried
func (t *Transport) shouldRetry(err error) bool {
	// Check error type to determine if retry is appropriate
	// e.g., network timeout, connection refused, etc.
	return true
}

// shouldRetryResponse determines if HTTP response should be retried
func (t *Transport) shouldRetryResponse(resp *Response) bool {
	// 5xx errors and 429 errors can be retried
	return resp.StatusCode >= 500 || resp.StatusCode == 429
}

// Close closes the transport layer
func (t *Transport) Close() error {
	// Close idle connections
	if transport, ok := t.client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
	return nil
}

// Helper function
func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}
