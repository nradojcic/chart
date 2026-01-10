package sitemap

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Remove trailing slash", "https://example.com/", "https://example.com"},
		{"Remove fragment", "https://example.com/page#section", "https://example.com/page"},
		{"Lower case host", "https://EXAMPLE.com/Page", "https://example.com/Page"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.expected {
				t.Errorf("got %s, want %s", got, tt.expected)
			}
		})
	}
}
