# Contributing to Go TabbyAPI Client

Thank you for your interest in contributing to the Go TabbyAPI client library! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)
- [Running Tests](#running-tests)
- [Benchmarks](#benchmarks)

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Getting Started

1. Fork the repository on GitHub
2. Clone your forked repository locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/go-tabbyapi.git
   cd go-tabbyapi
   ```

3. Add the original repository as an upstream remote:
   ```bash
   git remote add upstream https://github.com/pixelsquared/go-tabbyapi.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Install development tools:
   ```bash
   # Install golangci-lint (linter)
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.0
   ```

## Coding Standards

This project follows Go community standards and best practices:

### Code Formatting

- All code must be formatted with `gofmt`
- Use the standard Go conventions for naming (e.g., CamelCase for exported functions/variables, camelCase for unexported ones)
- Include comments for exported functions, types, and constants
- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines

### Tests

- All new functionality should include appropriate tests
- Maintain or improve code coverage with new contributions
- Tests should be contained in `*_test.go` files in the same package as the code being tested

### Documentation

- All exported types, functions, and methods must have proper documentation comments (godoc format)
- Update the README.md or other documentation for significant changes
- Include examples when adding new major functionality

## Pull Request Process

1. Ensure your code follows the coding standards
2. Update documentation as needed
3. Add or update tests for your changes
4. Run the test suite and ensure all tests pass
5. Make sure your branch is up to date with the main repository:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

6. Submit a pull request to the main repository's `main` branch
7. In your pull request description:
   - Clearly describe the problem you're solving
   - Include any relevant issue numbers using GitHub's #issue syntax
   - Explain your approach and any important design decisions

## Review Process

1. At least one project maintainer will review your PR
2. Automated CI checks must pass before merging
3. Address any feedback from reviewers
4. Once approved, a maintainer will merge your PR

## Release Process

### For Maintainers

1. Update CHANGELOG.md with new version information
2. Ensure all tests pass on the main branch
3. Create a new Git tag following semantic versioning:
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3"
   git push origin v1.2.3
   ```

4. Create a GitHub release with release notes
5. The new version will be available via Go modules

### Version Numbering

This project follows [Semantic Versioning](https://semver.org/):

- MAJOR version when making incompatible API changes
- MINOR version when adding functionality in a backward-compatible manner
- PATCH version when making backward-compatible bug fixes

## Running Tests

Run the full test suite:

```bash
make test
```

Or using Go directly:

```bash
go test -v -race -cover ./...
```

## Benchmarks

Run performance benchmarks:

```bash
make bench
```

Or using Go directly:

```bash
go test -run=^$ -bench=. -benchmem ./...
```

## Using the Makefile

This project includes a Makefile to simplify common development tasks:

- `make build`: Build the library
- `make test`: Run tests
- `make bench`: Run benchmarks
- `make lint`: Run linter
- `make fmt`: Format code
- `make tidy`: Tidy go.mod file
- `make examples`: Build examples
- `make check`: Run all checks (format, tidy, lint, test)
- `make help`: Show available commands

## License

By contributing to this project, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).

## Questions?

If you have any questions about contributing, please open an issue or contact the project maintainers.

Thank you for your contributions!