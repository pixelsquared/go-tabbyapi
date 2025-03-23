# Makefile for go-tabbyapi

# Configuration
GOBIN := $(shell go env GOPATH)/bin
GOLANGCI_LINT_VERSION := v1.56.0

# Default target
.PHONY: all
all: fmt lint test build examples

# Build the library
.PHONY: build
build:
	go build -v ./...

# Run tests
.PHONY: test
test:
	go test -v -race -cover ./...

# Run benchmarks
.PHONY: bench
bench:
	go test -run=^$ -bench=. -benchmem ./internal/rest ./internal/stream

# Run linter
.PHONY: lint
lint:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION); \
	fi
	$(GOBIN)/golangci-lint run ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Tidy go.mod file
.PHONY: tidy
tidy:
	go mod tidy

# Build examples
.PHONY: examples
examples:
	go build -v ./examples/...

# Clean build artifacts
.PHONY: clean
clean:
	go clean
	rm -rf bin/ dist/

# Generate documentation
.PHONY: docs
docs:
	@echo "Documentation generation not configured yet"

# Run all checks (useful before submitting PRs)
.PHONY: check
check: fmt tidy lint test

# Run checks without linting (useful when linter has compatibility issues)
.PHONY: test-only
test-only: fmt tidy test

# Help command
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all       : Run fmt, lint, test, build, and examples"
	@echo "  build     : Build the library"
	@echo "  test      : Run tests"
	@echo "  bench     : Run benchmarks"
	@echo "  lint      : Run linter"
	@echo "  fmt       : Format code"
	@echo "  tidy      : Tidy go.mod file"
	@echo "  examples  : Build examples"
	@echo "  clean     : Clean build artifacts"
	@echo "  docs      : Generate documentation"
	@echo "  check     : Run all checks (fmt, tidy, lint, test)"
	@echo "  test-only : Run checks without linting"
	@echo "  help      : Show this help message"