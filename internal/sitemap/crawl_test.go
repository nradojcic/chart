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

	found := false
	want := server.URL + "/page1"
	for _, p := range pages {
		if p == want {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find %s in results, but it was missing. Got: %v", want, pages)
	}
}

func TestCrawlRelativeLinks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if r.URL.Path == "/" {
			w.Write([]byte(`
				<html>
					<body>
						<a href="projects/chart">Relative Link</a>
						<a href="/contact">Root Relative Link</a>
					</body>
				</html>
			`))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	ctx := context.Background()
	pages := Crawl(ctx, server.URL, 1, "test-agent", 1, nil)

	expected := []string{
		server.URL,
		server.URL + "/projects/chart",
		server.URL + "/contact",
	}

	if len(pages) != len(expected) {
		t.Errorf("Expected %d pages, got %d", len(expected), len(pages))
	}

	for _, want := range expected {
		found := false
		for _, p := range pages {
			if p == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find %s, but it was missing", want)
		}
	}
}
