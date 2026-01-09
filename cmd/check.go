package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

type CheckResult struct {
	Url    string
	Status string
	Code   int
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks a list of URLs from stdin to see if they are live or dead.",
	RunE: func(cmd *cobra.Command, args []string) error {
		resultsChan := make(chan CheckResult)
		var wg sync.WaitGroup

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			url := scanner.Text()
			wg.Add(1)
			go checkUrl(url, resultsChan, &wg)
		}

		go func() {
			wg.Wait()
			close(resultsChan)
		}()

		var liveLinks, deadLinks []CheckResult
		for result := range resultsChan {
			if result.Status == "live" {
				liveLinks = append(liveLinks, result)
			} else {
				deadLinks = append(deadLinks, result)
			}
		}

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkUrl(url string, resultsChan chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		resultsChan <- CheckResult{Url: url, Status: "dead", Code: 0}
		return
	}
	req.Header.Set("User-Agent", UserAgent)

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
