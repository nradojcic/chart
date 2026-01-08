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
	Run: func(cmd *cobra.Command, args []string) {
		urlStr := args[0]
		maxDepth := viper.GetInt("depth")
		outputFormat := viper.GetString("format")

		pages := sitemap.Crawl(urlStr, maxDepth)

		// Text output
		if outputFormat == "txt" {
			for _, page := range pages {
				fmt.Println(page)
			}
			return
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
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().IntP("depth", "d", 3, "The maximum depth to traverse")
	buildCmd.Flags().StringP("format", "f", "xml", "The output format (xml or txt)")

	viper.BindPFlag("depth", buildCmd.Flags().Lookup("depth"))
	viper.BindPFlag("format", buildCmd.Flags().Lookup("format"))
}
