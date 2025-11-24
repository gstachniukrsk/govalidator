.PHONY: help test coverage coverage-html lint fmt clean

# Default target
help:
	@echo "Available targets:"
	@echo "  make test          - Run all tests"
	@echo "  make coverage      - Run tests with coverage (requires 95%)"
	@echo "  make coverage-html - Generate HTML coverage report and open in browser"
	@echo "  make lint          - Run linters"
	@echo "  make fmt           - Format code"
	@echo "  make clean         - Clean generated files"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run coverage check (requires 95%)
coverage:
	@echo "Running coverage check (requires 95%)..."
	@./scripts/coverage.sh

# Generate HTML coverage report and open in browser
coverage-html: coverage
	@echo "Opening coverage report in browser..."
	@if command -v open > /dev/null; then \
		open coverage.html; \
	elif command -v xdg-open > /dev/null; then \
		xdg-open coverage.html; \
	else \
		echo "Please open coverage.html manually"; \
	fi

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -f coverage.out coverage.html
	go clean -testcache
