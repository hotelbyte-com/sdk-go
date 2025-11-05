package hotelbyte

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

// Transport represents HTTP transport layer
type Transport struct {
	client *resty.Client
	config *Config
}

// NewTransport creates a new transport layer
func NewTransport(config *Config) (*Transport, error) {
	client := resty.New().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.HTTPConfig.Timeout).
		SetHeader("User-Agent", config.HTTPConfig.UserAgent).
		SetHeader("Content-Type", "application/json").
		SetJSONMarshaler(sonic.Marshal).
		SetJSONUnmarshaler(sonic.Unmarshal).
		SetTransport(&http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			MaxIdleConns:        config.HTTPConfig.MaxIdleConns,
			MaxIdleConnsPerHost: config.HTTPConfig.MaxConnsPerHost,
			IdleConnTimeout:     90 * time.Second,
		}).
		SetRetryCount(config.RetryConfig.MaxRetries).
		SetRetryWaitTime(config.RetryConfig.InitialDelay).
		SetRetryMaxWaitTime(config.RetryConfig.MaxDelay).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// 重试条件：网络错误 || 429 Too Many Requests || 5xx Server Error
			return err != nil || r.StatusCode() == 429 || r.StatusCode() >= 500
		})

	return &Transport{
		client: client,
		config: config,
	}, nil
}

// Do executes HTTP request
func (t *Transport) Do(ctx context.Context, req *types.HttpRequest) (*types.HttpResponse, error) {
	// Build Resty request
	r := t.client.R().SetContext(ctx)

	// Set custom headers
	if req.Headers != nil {
		r.SetHeaders(req.Headers)
	}

	// Set query parameters
	if req.Query != nil {
		r.SetQueryParamsFromValues(req.Query)
	}

	// Set request body
	if req.Body != nil {
		r.SetBody(req.Body)
	}

	// Execute request (Resty handles retry internally)
	resp, err := r.Execute(req.Method, req.Path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	headers := resp.Header()
	sb := strings.Builder{}
	for _, key := range keys {
		if val := headers.Get(key); val != "" {
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(val)
			sb.WriteString(" ")
		}
	}
	if sb.Len() > 0 {
		logrus.WithContext(ctx).Infof("%s Response headers: %s", req.Path, strings.TrimSpace(sb.String()))
	}
	return &types.HttpResponse{
		StatusCode: resp.StatusCode(),
		Headers:    resp.Header(),
		Body:       resp.Body(),
	}, nil
}

var (
	keys = []string{
		"Trace-Id",
		"Session-Id",
		"Server-Cost-Milliseconds",
	}
)

// Close closes the transport layer
func (t *Transport) Close() error {
	// Close idle connections
	if transport, ok := t.client.GetClient().Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
	return nil
}
