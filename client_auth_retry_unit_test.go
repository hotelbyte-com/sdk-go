package hotelbyte

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hotelbyte-com/sdk-go/protocol/types"
)

func TestDoWithAuthRetry_SkipRefreshWhenTokenAlreadyChanged(t *testing.T) {
	var authCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth/ticket" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		atomic.AddInt32(&authCalls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"data":{"ticket":"refreshed-by-server"}}`))
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithCredentials("app", "secret"),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	client.SetToken("old-token", 3600)
	ctx := context.Background()

	attempts := 0
	result, err := doWithAuthRetry[int](ctx, client, func() (*int, error) {
		attempts++
		if attempts == 1 {
			// Simulate another in-flight request that already refreshed token.
			client.SetToken("new-token", 3600)
			return nil, types.NewBizErr(ErrCodeTokenExpired, "token expired")
		}
		v := 1
		return &v, nil
	})
	if err != nil {
		t.Fatalf("doWithAuthRetry error: %v", err)
	}
	if result == nil || *result != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
	if got := atomic.LoadInt32(&authCalls); got != 0 {
		t.Fatalf("refresh auth calls = %d, want 0", got)
	}
}

func TestDoWithAuthRetry_RefreshWhenTokenUnchanged(t *testing.T) {
	var authCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth/ticket" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		atomic.AddInt32(&authCalls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"data":{"ticket":"fresh-token"}}`))
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithCredentials("app", "secret"),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	client.SetToken("old-token", 3600)
	ctx := context.Background()

	attempts := 0
	result, err := doWithAuthRetry[int](ctx, client, func() (*int, error) {
		attempts++
		if attempts == 1 {
			return nil, types.NewBizErr(ErrCodeTokenExpired, "token expired")
		}
		if got := client.GetToken(); got != "fresh-token" {
			t.Fatalf("token after refresh = %q, want fresh-token", got)
		}
		v := 2
		return &v, nil
	})
	if err != nil {
		t.Fatalf("doWithAuthRetry error: %v", err)
	}
	if result == nil || *result != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
	if got := atomic.LoadInt32(&authCalls); got != 1 {
		t.Fatalf("refresh auth calls = %d, want 1", got)
	}
}

func TestDoWithAuthRetry_DisabledBypassesAuthAndRetry(t *testing.T) {
	var authCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth/ticket" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		atomic.AddInt32(&authCalls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"data":{"ticket":"fresh-token"}}`))
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithCredentials("app", "secret"),
		WithDisableAutoAuthRetry(true),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	ctx := context.Background()
	attempts := 0
	_, err = doWithAuthRetry[int](ctx, client, func() (*int, error) {
		attempts++
		return nil, types.NewBizErr(ErrCodeTokenExpired, "token expired")
	})
	if err == nil {
		t.Fatalf("expected error when auto auth retry disabled")
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
	if got := atomic.LoadInt32(&authCalls); got != 0 {
		t.Fatalf("refresh auth calls = %d, want 0", got)
	}
}
