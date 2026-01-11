package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestCheckUrl(t *testing.T) {
	// Create a mock server to simulate live and dead pages
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/live" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	ctx := context.Background()
	resultsChan := make(chan CheckResult, 2)
	var wg sync.WaitGroup
	guard := make(chan struct{}, 1)

	// Test 1: Check a URL that should be "live"
	wg.Add(1)
	guard <- struct{}{}
	CheckUrl(ctx, server.URL+"/live", resultsChan, &wg, guard, "test-bot", nil)

	got := <-resultsChan
	if got.Status != "live" || got.Code != 200 {
		t.Errorf("Expected live/200, got %s/%d", got.Status, got.Code)
	}

	// Test 2: Check a URL that should be "dead"
	wg.Add(1)
	guard <- struct{}{}
	CheckUrl(ctx, server.URL+"/404", resultsChan, &wg, guard, "test-bot", nil)

	got = <-resultsChan
	if got.Status != "dead" || got.Code != 404 {
		t.Errorf("Expected dead/404, got %s/%d", got.Status, got.Code)
	}

	// Test 3: Check an invalid URL
	wg.Add(1)
	guard <- struct{}{}
	CheckUrl(ctx, "http://invalid-url:", resultsChan, &wg, guard, "test-bot", nil)

	got = <-resultsChan
	if got.Status != "dead" || got.Code != 0 {
		t.Errorf("Expected dead/0, got %s/%d", got.Status, got.Code)
	}
}
