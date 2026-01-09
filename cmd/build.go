package cmd

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/nradojcic/chart/internal/sitemap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [url]",
	Short: "Builds a sitemap for the provided URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		urlStr := args[0]
		maxDepth := viper.GetInt("depth")
		outputFormat := viper.GetString("format")
		userAgent := viper.GetString("user-agent")
		concurrency := viper.GetInt("concurrency")

		const maxConcurrency = 100 // upper limit on user provided input to avoid resource exhaustion
		if concurrency > maxConcurrency {
			concurrency = maxConcurrency
		}
		if concurrency < 1 {
			concurrency = 1
		}

		pages := sitemap.Crawl(urlStr, maxDepth, userAgent, concurrency)

		// Text output
		if outputFormat == "txt" {
			for _, page := range pages {
				fmt.Println(page)
			}
			return nil
		}

		// Default XML output
		toXml := urlset{
			Xmlns: xmlns,
		}
		for _, page := range pages {
			toXml.Urls = append(toXml.Urls, loc{page})
		}

		fmt.Print(xml.Header)
		enc := xml.NewEncoder(os.Stdout)
		enc.Indent("", "  ")
		if err := enc.Encode(toXml); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding XML: %v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().IntP("depth", "d", 3, "The maximum depth to traverse")
	buildCmd.Flags().StringP("format", "f", "xml", "The output format (xml or txt)")

	viper.BindPFlag("depth", buildCmd.Flags().Lookup("depth"))
	viper.BindPFlag("format", buildCmd.Flags().Lookup("format"))
}
