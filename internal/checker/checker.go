package checker

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type CheckResult struct {
	Url    string
	Status string
	Code   int
}

// CheckUrl performs a HEAD request to the given URL and sends the result to resultsChan
func CheckUrl(ctx context.Context, url string, resultsChan chan<- CheckResult, wg *sync.WaitGroup, guard chan struct{}, userAgent string, throttle <-chan time.Time) {
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

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		resultsChan <- CheckResult{Url: url, Status: "dead", Code: 0}
		return
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		resultsChan <- CheckResult{Url: url, Status: "dead", Code: 0}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		resultsChan <- CheckResult{Url: url, Status: "live", Code: resp.StatusCode}
	} else {
		resultsChan <- CheckResult{Url: url, Status: "dead", Code: resp.StatusCode}
	}
}
