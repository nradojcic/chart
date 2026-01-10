package link

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		html string
		want int
	}{
		{"Single link", `<html><body><a href="/1">Link 1</a></body></html>`, 1},
		{"Multiple links", `<html><body><a href="/1">1</a><a href="/2">2</a></body></html>`, 2},
		{"No links", `<html><body><p>Hello World</p></body></html>`, 0},
		{"Nested tags", `<html><body><a href="/3"><span>Text</span></a></body></html>`, 1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.html)
			links, err := Parse(r)
			if err != nil {
				t.Fatalf("Parse() failed: %v", err)
			}
			if len(links) != tt.want {
				t.Errorf("expected %d links, got %d", tt.want, len(links))
			}
		})
	}
}
