#!/bin/bash
set -e

# Coverage check script
# Requires 95% code coverage to pass

COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"
MIN_COVERAGE=95.0

echo "Running tests with coverage..."
go test -v -coverprofile="$COVERAGE_FILE" -covermode=atomic ./...

echo ""
echo "Generating coverage report..."
go tool cover -func="$COVERAGE_FILE"

echo ""
echo "Calculating total coverage..."
TOTAL_COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' | sed 's/%//')

echo ""
echo "=========================================="
echo "Total Coverage: ${TOTAL_COVERAGE}%"
echo "Required Coverage: ${MIN_COVERAGE}%"
echo "=========================================="

# Generate HTML report
go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
echo ""
echo "HTML coverage report generated: $COVERAGE_HTML"

# Check if coverage meets minimum threshold
if (( $(echo "$TOTAL_COVERAGE < $MIN_COVERAGE" | bc -l) )); then
    echo ""
    echo "❌ FAILED: Coverage ${TOTAL_COVERAGE}% is below required ${MIN_COVERAGE}%"
    echo ""
    echo "Uncovered lines:"
    go tool cover -func="$COVERAGE_FILE" | grep -v "100.0%" | grep -v "total:"
    exit 1
else
    echo ""
    echo "✅ PASSED: Coverage meets requirement of ${MIN_COVERAGE}%"
    exit 0
fi
