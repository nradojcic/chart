package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "chart",
	Short:   "chart is a fast and flexible sitemap generator",
	Long: `Site Chart (chart) is a CLI tool designed to crawl a website 
and generate a comprehensive sitemap. 

It follows links recursively within the same domain to 
map out your site's structure. You can customize the depth of the 
crawl and the output format via flags or a configuration file.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chart.yaml)")
	rootCmd.PersistentFlags().IntP("concurrency", "c", 10, "Number of concurrent workers")
	rootCmd.PersistentFlags().Float64P("rate-limit", "r", 2, "Rate limit in requests per second (0 for no limit)")
	rootCmd.PersistentFlags().StringP("user-agent", "u", "SiteChart-Sitemapper/1.0", "Custom User-Agent header for HTTP requests")
	rootCmd.PersistentFlags().DurationP("timeout", "t", 0, "Global timeout for the command (e.g., 30s, 1m, 1h)")

	viper.BindPFlag("concurrency", rootCmd.PersistentFlags().Lookup("concurrency"))
	viper.BindPFlag("rate-limit", rootCmd.PersistentFlags().Lookup("rate-limit"))
	viper.BindPFlag("user-agent", rootCmd.PersistentFlags().Lookup("user-agent"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".chart" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".chart")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
