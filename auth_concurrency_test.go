package hotelbyte

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_ConcurrentRefreshToken(t *testing.T) {
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/auth/ticket" {
			atomic.AddInt32(&requestCount, 1)
			time.Sleep(10 * time.Millisecond) // Simulate network delay
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code":0,"data":{"ticket":"ST:mock-token-` + time.Now().String() + `"}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithCredentials("test-app", "test-secret"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Simulate 100 concurrent RefreshToken calls
	var wg sync.WaitGroup
	numRoutines := 100
	wg.Add(numRoutines)

	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			err := client.RefreshToken(ctx)
			if err != nil {
				t.Errorf("RefreshToken failed: %v", err)
			}
		}()
	}

	wg.Wait()

	// Thanks to singleflight, the number of actual requests to the server should be 1
	if count := atomic.LoadInt32(&requestCount); count != 1 {
		t.Errorf("Expected 1 request to /api/auth/ticket, got %d", count)
	}
}
