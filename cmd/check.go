package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/nradojcic/chart/internal/checker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [urls...]",
	Short: "Checks if provided URLs are live or dead",
	Long:  `Provide URLs as arguments or pipe them via stdin.`,
	Example: `  chart check https://example1.com https://example2.com
  cat urls.txt | chart check`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		timeout := viper.GetDuration("timeout")
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		resultsChan := make(chan checker.CheckResult)
		var wg sync.WaitGroup
		var urls []string

		userAgent := viper.GetString("user-agent")
		concurrency := viper.GetInt("concurrency")
		rateLimit := viper.GetFloat64("rate-limit")

		// Gather input from Args or Stdin
		if len(args) > 0 {
			urls = args
		} else {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				return fmt.Errorf("no URLs provided via arguments or stdin")
			}
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				urls = append(urls, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return err
			}
		}

		// Check URLs
		const maxConcurrency = 100 // upper limit on user provided concurrency to avoid resource exhaustion
		if concurrency > maxConcurrency {
			concurrency = maxConcurrency
		}
		if concurrency < 1 {
			concurrency = 1
		}

		const maxRateLimit = 100.0 // upper limit on user provided rate to avoid abuse
		if rateLimit > maxRateLimit {
			rateLimit = maxRateLimit
		}
		if rateLimit < 0 {
			rateLimit = 0
		}

		var throttle <-chan time.Time
		if rateLimit > 0 {
			ticker := time.NewTicker(time.Duration(float64(time.Second) / rateLimit))
			defer ticker.Stop()
			throttle = ticker.C
		}

		guard := make(chan struct{}, concurrency) // semaphore to limit concurrency

	loop:
		for _, url := range urls {
			select {
			case <-ctx.Done():
				break loop
			case guard <- struct{}{}: // block when guard channel capacity full
				wg.Add(1)
				go checker.CheckUrl(ctx, url, resultsChan, &wg, guard, userAgent, throttle)
			}
		}

		go func() {
			wg.Wait()
			close(resultsChan)
		}()

		// Process results
		var liveLinks, deadLinks []checker.CheckResult
		for result := range resultsChan {
			if result.Status == "live" {
				liveLinks = append(liveLinks, result)
			} else {
				deadLinks = append(deadLinks, result)
			}
		}

		sort.Slice(liveLinks, func(i, j int) bool {
			return liveLinks[i].Url < liveLinks[j].Url
		})
		sort.Slice(deadLinks, func(i, j int) bool {
			return deadLinks[i].Url < deadLinks[j].Url
		})

		fmt.Println("Live Links:")
		for _, link := range liveLinks {
			fmt.Printf("  [%d] %s\n", link.Code, link.Url)
		}

		fmt.Println("\nDead Links:")
		for _, link := range deadLinks {
			fmt.Printf("  [%d] %s\n", link.Code, link.Url)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
