package sitemap

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Remove trailing slash", "https://example.com/", "https://example.com"},
		{"Remove fragment", "https://example.com/page#section", "https://example.com/page"},
		{"Lower case host", "https://EXAMPLE.com/Page", "https://example.com/Page"},
		{"Handle empty string", "", ""},
		{"Already normalized", "https://example.com/test", "https://example.com/test"},
		{"Keep query parameters", "https://example.com?q=1", "https://example.com?q=1"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}
