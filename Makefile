# Variables
BINARY_NAME=chart
BIN_DIR=bin

.PHONY: all build test coverage clean help

all: test build

## build: Compiles the binary into the bin directory
build:
	@echo "Building..."
	go build -o $(BIN_DIR)/$(BINARY_NAME) main.go

## test: Runs all tests with the race detector enabled
test:
	@echo "Running tests..."
	go test -v -race ./...

## coverage: Runs tests and generates a coverage report
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## clean: Removes the compiled binary and build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html
	go clean

## help: Shows this help message
help:
	@echo "Usage: make [target]"
	@grep -E '^##' Makefile | sed -e 's/## //g' | column -t -s ':' | sed -e 's/^/  /'
