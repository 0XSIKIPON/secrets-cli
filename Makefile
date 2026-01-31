.PHONY: build clean test test-e2e install build-all

# Binary name
BINARY=secrets-cli

# Build info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "local-build")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target
all: build

# Build the binary
build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/secrets-cli

# Build for all platforms (for releases)
build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 ./cmd/secrets-cli
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 ./cmd/secrets-cli
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 ./cmd/secrets-cli
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 ./cmd/secrets-cli

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/secrets-cli

# Clean build artifacts
clean:
	rm -f $(BINARY)
	rm -rf dist/
	rm -f coverage.out

# Run Go tests
test:
	go test -v ./...

# Run E2E tests (requires Docker)
test-e2e: build
	./tests/run-tests.sh --rebuild

# Run E2E tests with verbose output
test-e2e-verbose: build
	./tests/run-tests.sh --rebuild -v

# Check code formatting
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Show help
help:
	@echo "Usage: make [target] [VERSION=vX.Y.Z]"
	@echo ""
	@echo "Targets:"
	@echo "  build           - Build the secrets-cli binary"
	@echo "  build-all       - Build for all platforms (linux/darwin, amd64/arm64)"
	@echo "  install         - Install to GOPATH/bin"
	@echo "  clean           - Remove build artifacts"
	@echo "  test            - Run Go unit tests"
	@echo "  test-e2e        - Run E2E tests in Docker"
	@echo "  fmt             - Format Go code"
	@echo "  lint            - Run golangci-lint"
	@echo ""
	@echo "Examples:"
	@echo "  make build                  # Build with auto-detected version"
	@echo "  make build VERSION=v1.0.0   # Build with specific version"

