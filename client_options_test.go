package hotelbyte

import (
	"testing"
	"time"
)

func TestWithHTTPConfig_PreservesHeadersSetByWithHeader(t *testing.T) {
	cfg := DefaultConfig()

	if err := WithHeader("X-Load-Test", "1")(cfg); err != nil {
		t.Fatalf("set header: %v", err)
	}

	httpCfg := HTTPConfig{
		Timeout:         3 * time.Second,
		MaxIdleConns:    10,
		MaxConnsPerHost: 20,
		UserAgent:       "test-agent",
		DisableHTTP2:    true,
	}
	if err := WithHTTPConfig(httpCfg)(cfg); err != nil {
		t.Fatalf("set http config: %v", err)
	}

	if got := cfg.HTTPConfig.DefaultHeaders["X-Load-Test"]; got != "1" {
		t.Fatalf("expected X-Load-Test header to be preserved, got %q", got)
	}
}

func TestWithHTTPConfig_HeaderInConfigTakesPrecedence(t *testing.T) {
	cfg := DefaultConfig()

	if err := WithHeader("X-Load-Test", "1")(cfg); err != nil {
		t.Fatalf("set header: %v", err)
	}

	httpCfg := HTTPConfig{
		Timeout:         3 * time.Second,
		MaxIdleConns:    10,
		MaxConnsPerHost: 20,
		UserAgent:       "test-agent",
		DefaultHeaders: map[string]string{
			"X-Load-Test": "true",
			"X-Custom":    "v",
		},
	}
	if err := WithHTTPConfig(httpCfg)(cfg); err != nil {
		t.Fatalf("set http config: %v", err)
	}

	if got := cfg.HTTPConfig.DefaultHeaders["X-Load-Test"]; got != "true" {
		t.Fatalf("expected http config header to win, got %q", got)
	}
	if got := cfg.HTTPConfig.DefaultHeaders["X-Custom"]; got != "v" {
		t.Fatalf("expected custom header to exist, got %q", got)
	}
}
