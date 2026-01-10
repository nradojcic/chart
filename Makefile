# Variables
BINARY_NAME=chart
BIN_DIR=bin

.PHONY: all build test clean help

all: test build

## build: Compiles the binary into the bin directory
build:
	@echo "Building..."
	go build -o $(BIN_DIR)/$(BINARY_NAME) main.go

## test: Runs all unit tests with the race detector enabled
test:
	@echo "Running tests..."
	go test -v -race ./...

## clean: Removes the compiled binary and build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)
	go clean

## help: Shows this help message
help:
	@echo "Usage: make [target]"
	@sed -n 's/^##//p' $(MAKEFILE) | column -t -s ':' |  sed -e 's/^/ /'
