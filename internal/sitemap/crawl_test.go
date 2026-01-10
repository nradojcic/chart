package sitemap

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrawl(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body><a href="/page1">Link</a></body></html>`))
	}))
	defer server.Close()

	// Run crawler against the mock server
	ctx := context.Background()
	pages := Crawl(ctx, server.URL, 1, "test-agent", 1, nil)

	if len(pages) == 0 {
		t.Fatal("Expected to find pages, but found none")
	}
}
