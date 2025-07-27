# Makefile for murl - HTTP benchmarking tool

# Variables
BINARY_NAME=murl
BUILD_DIR=build
CMD_DIR=cmd/murl
MAIN_PACKAGE=./$(CMD_DIR)
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Default target
.PHONY: all
all: clean build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux amd64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	# Linux arm64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	# macOS amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	# macOS arm64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	# Windows amd64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "Cross-compilation complete"

# Install the binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Build and run test server (current platform)
.PHONY: test-server
test-server:
	@echo "Building and starting test server..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/testserver ./cmd/testserver
	./$(BUILD_DIR)/testserver

# Build test server for multiple platforms
.PHONY: test-server-all
test-server-all: clean
	@echo "Building test server for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux amd64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/testserver-linux-amd64 ./cmd/testserver
	# Linux arm64
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/testserver-linux-arm64 ./cmd/testserver
	# macOS amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/testserver-darwin-amd64 ./cmd/testserver
	# macOS arm64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/testserver-darwin-arm64 ./cmd/testserver
	@echo "Cross-platform test server binaries built successfully!"
	@ls -la $(BUILD_DIR)/testserver-*

# Build and run test server for Linux (useful for testing in containers)
.PHONY: test-server-linux
test-server-linux:
	@echo "Building test server for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/testserver-linux-amd64 ./cmd/testserver
	@echo "Linux test server built: $(BUILD_DIR)/testserver-linux-amd64"

# Build and run test server for macOS
.PHONY: test-server-darwin
test-server-darwin:
	@echo "Building test server for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/testserver-darwin-amd64 ./cmd/testserver
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/testserver-darwin-arm64 ./cmd/testserver
	@echo "macOS test server built: $(BUILD_DIR)/testserver-darwin-*"

# Run unit tests only
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v ./internal/...

# Run integration tests
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v ./test/...

# Run all tests (unit + integration)
.PHONY: test-all
test-all: test-unit test-integration
	@echo "All tests completed"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run the application (for development)
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) --help

# Development setup
.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	$(GOMOD) download
	@echo "Installing development tools..."
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development setup complete"

# Quick development build and test
.PHONY: dev
dev: fmt build test
	@echo "Development build and test complete"

# Release build (optimized)
.PHONY: release
release: clean
	@echo "Building release version..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Release build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Docker build (if needed in the future)
.PHONY: docker
docker:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build (default)"
	@echo "  build        - Build the binary"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  bench        - Run benchmarks"
	@echo "  test-server  - Build and run test server (current platform)"
	@echo "  test-server-all - Build test server for multiple platforms"
	@echo "  test-server-linux - Build test server for Linux"
	@echo "  test-server-darwin - Build test server for macOS"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  tidy         - Tidy dependencies"
	@echo "  deps         - Download dependencies"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run with --help"
	@echo "  dev-setup    - Setup development environment"
	@echo "  dev          - Quick development build and test"
	@echo "  release      - Build optimized release version"
	@echo "  docker       - Build Docker image"
	@echo "  help         - Show this help"
