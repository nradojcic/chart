package link

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	html := `
	<html>
	<body>
		<h1>My first page</h1>
		<p>This is my first page</p>
		<a href="/page2">Second page</a>
	</body>
	</html>
	`

	r := strings.NewReader(html)
	links, err := Parse(r)
	if err != nil {
		t.Fatalf("Parse() returned an error: %v", err)
	}

	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}

	if links[0].Href != "/page2" {
		t.Errorf("expected href '/page2', got '%s'", links[0].Href)
	}

	if links[0].Text != "Second page" {
		t.Errorf("expected text 'Second page', got '%s'", links[0].Text)
	}
}
