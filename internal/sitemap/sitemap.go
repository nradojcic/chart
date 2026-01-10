package sitemap

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nradojcic/chart/internal/link"
)

// Crawl performs a breadth-first crawl of the website starting at urlStr
// up to maxDepth levels deep, returning a slice of discovered URLs.
func Crawl(ctx context.Context, urlStr string, maxDepth int, userAgent string, concurrency int, throttle <-chan time.Time) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		Normalize(urlStr): {},
	}

	for i := 0; i <= maxDepth; i++ {
		// check for context cancellation at the beginning of each level of crawl
		select {
		case <-ctx.Done():
			break
		default:
		}

		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}

		var urlsToCrawl []string
		for url := range q {
			if _, ok := seen[url]; !ok {
				urlsToCrawl = append(urlsToCrawl, url)
				seen[url] = struct{}{}
			}
		}

		if len(urlsToCrawl) == 0 {
			continue
		}

		linksChan := make(chan []string)
		var wg sync.WaitGroup
		guard := make(chan struct{}, concurrency) // semaphore to limit concurrency

	loop:
		for _, url := range urlsToCrawl {
			select {
			case <-ctx.Done():
				break loop
			case guard <- struct{}{}: // block when guard channel capacity full
				wg.Add(1)
				go func(u string) {
					defer func() {
						wg.Done()
						<-guard
					}()

					if throttle != nil {
						select {
						case <-throttle:
						case <-ctx.Done():
							return
						}
					}
					linksChan <- get(ctx, u, userAgent)
				}(url)
			}
		}

		go func() {
			wg.Wait()
			close(linksChan)
		}()

		for links := range linksChan {
			for _, link := range links {
				normalizedLink := Normalize(link)
				if _, ok := seen[normalizedLink]; !ok {
					nq[normalizedLink] = struct{}{}
				}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}

	sort.Strings(ret)
	return ret
}

func get(ctx context.Context, urlStr string, userAgent string) []string {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return []string{}
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	return filter(hrefs(resp.Body, base), withPrefix(base)) // only keep links from the same base domain
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}

	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

// Normalize normalizes the URL by removing fragments and trailing slashes
func Normalize(rawUrl string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	u.Fragment = ""                            // remove # fragment from URL
	return strings.TrimSuffix(u.String(), "/") // remove trailing slash
}
