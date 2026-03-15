# Project-specific justfile
# https://github.com/casey/just

# Default recipe - show available commands
default:
    @just --list

# Run all tests
test:
    go test -v -race -cover ./...

# Run tests with coverage report
test-coverage:
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report: coverage.html"

# Run tests with verbose output
test-verbose:
    go test -v -race ./...

# Run example tests
test-examples:
    go test -v -run Example ./...

# Run benchmarks
bench:
    go test -bench=. -benchmem ./...

# Run fuzzing tests (requires Go 1.18+)
fuzz:
    go test -fuzz=Fuzz -fuzztime=30s ./...

# Run linter
lint:
    golangci-lint run ./...

# Run linter with auto-fix
lint-fix:
    golangci-lint run --fix ./...

# Run go vet
vet:
    go vet ./...

# Run security scanner
security:
    gosec ./...

# Build the package
build:
    go build ./...

# Clean build artifacts and caches
clean:
    go clean -cache -testcache -modcache -i -r

# Download dependencies
deps:
    go mod download
    go mod tidy

# Update dependencies
deps-update:
    go get -u ./...
    go mod tidy

# Show outdated dependencies
deps-outdated:
    go list -u -m -json all | go-mod-outdated -update -direct

# Generate godoc documentation
doc:
    godoc -http=:6060

# Format code
fmt:
    go fmt ./...

# Check formatting
fmt-check:
    @gofmt -l . | grep -q . && echo "Code needs formatting" && exit 1 || echo "Code is formatted"

# Run all quality checks
check: fmt-check vet lint test
    echo "All checks passed!"

# Run pre-commit checks
pre-commit: fmt vet test
    echo "Pre-commit checks passed!"

# Install pre-commit hooks
install-hooks:
    cp .git/hooks/pre-commit.sample .git/hooks/pre-commit 2>/dev/null || true
    echo '#!/bin/sh\njust pre-commit' > .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    echo "Pre-commit hooks installed"

# Show project info
info:
    @echo "Project: github.com/artmann/businessrules"
    @echo "Go version: $$(go version)"
    @echo "Module info:"
    @go list -m all | head -10

# Tag a new release
tag version:
    git tag -a v{{version}} -m "Release v{{version}}"
    git push origin v{{version}}
    @echo "Tagged and pushed v{{version}}"

# Bump version in doc.go
bump-version version:
    sed -i 's/Version = ".*"/Version = "{{version}}"/' doc.go
    @echo "Updated version to {{version}}"

# Dogfood - use this library's patterns on itself
dogfood:
    go test -v -race -cover ./...
    go vet ./...
    go build ./...
    @echo "Dogfooding complete - library is healthy!"
