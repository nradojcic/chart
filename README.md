# Site Chart (`chart`)

[![Go Version](https://img.shields.io/github/go-mod/go-version/nradojcic/chart)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Site Chart** is a high-performance CLI tool built in Go for web developers and SEO specialists. It provides two primary functions: recursively crawling a website to generate a comprehensive sitemap and checking a list of URLs for availability (live/dead links).

Designed with efficiency in mind, `chart` leverages Go's concurrency to handle large-scale crawls and link checks swiftly while respecting server limits.

## Technical Highlights

- **Concurrency Control**: Implemented a worker-pool pattern using buffered channels (semaphores) to manage concurrent HTTP requests efficiently without exhausting system resources.
- **Rate Limiting**: Utilized a custom ticker-based throttler to ensure compliance with web server policies and prevent IP rate-limiting.
- **Robust Testing**: Achieved 80%+ test coverage in core packages using table-driven unit tests and httptest for hermetic network mocking.
- **Graceful Shutdown**: Context-aware architecture allows for clean cancellation of long-running crawls via OS signals (Ctrl+C), or global timeout, ensuring all goroutines exit cleanly.

##  Key Features

- **Recursive Sitemap Generation**: Map out entire site structure with configurable crawling depth.
- **Concurrent Execution**: High-speed processing using Go routines and worker pools.
- **Rate Limiting**: Built-in throttling to ensure you don't overwhelm target servers.
- **Link Validator**: Batch check URLs for status codes and availability via CLI arguments or `stdin`.
- **Flexible Configuration**: Manage settings via command-line flags, environment variables, or YAML config files.
- **Output Formats**: Export sitemaps in standard XML or plain text formats.
- **Graceful Termination**: Handles OS signals or global timeout for safe shutdowns during long-running tasks.

##  Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (1.25 or later recommended)
- Make (for building and testing)

### Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/nradojcic/chart.git
cd chart
make build
```

The executable will be available in the local `bin/` directory. You can add it to your PATH:

```bash
cp bin/chart /usr/local/bin/
```

##  Usage

`chart` provides a clean, intuitive CLI interface powered by Cobra.

### 1. Generating a Sitemap

Crawl a website up to a specific depth and output in XML (default):

```bash
# Crawl a website and print a TXT sitemap to standard output
chart build https://example.com --format=txt

# Crawl a website up to a depth of 5 and save as an XML sitemap file
chart build https://example.com --depth 5 > sitemap.xml
```

### 2. Checking Link Status

Check live vs dead links from a list of URLs:

```bash
# Check individual URLs
chart check https://example1.com https://example2.com

# Check URLs from a local file
cat urls.txt | chart check

# Chain commands: build a sitemap and immediately check it for broken links
chart build https://example.com --format=txt | chart check
```

### Global Flags

- `-c, --concurrency int`: Number of concurrent workers (default 10).
- `-r, --rate-limit float`: Max requests per second (default 2).
- `-u, --user-agent string`: Custom User-Agent header for identification to websites.
- `-t, --timeout duration`: Global timeout (e.g., `30s`, `5m`).

## Configuration

`chart` uses Viper for flexible configuration. It looks for a `.chart.yaml` file in the current directory or your home folder.

Example `.chart.yaml`:

```yaml
concurrency: 20              # number of concurrent workers
rate-limit: 1.5              # requests per second
user-agent: "ChartBot/1.0"   # custom bot identification
depth: 5                     # maximum crawl depth
format: "xml"                # output format: xml or txt
timeout: 30m                 # global timeout duration
```

## Tech Stack

- **[Go](https://golang.org/)**: The core language, chosen for its performance and concurrency model.
- **[Cobra](https://github.com/spf13/cobra)**: CLI framework for modern command-line interfaces.
- **[Viper](https://github.com/spf13/viper)**: Configuration management for Go applications.
- **Standard Library**: Extensively used for networking (`net/http`), XML processing, and concurrency.

## Development

Run the test suite with the race detector enabled:

```bash
make test
```

To view the visual line-by-line coverage report `coverage.html`:
```bash
make coverage
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
