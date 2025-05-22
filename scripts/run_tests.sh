#!/bin/bash
# ==========================================
# run_tests.sh - Run tests and generate coverage report
# ==========================================
set -e  # Exit on error

# Run tests with coverage
echo "Running tests with coverage..."
go test -v -coverprofile=coverage.txt ./...

# Display coverage report
go tool cover -func=coverage.txt

echo "Coverage report saved to coverage.txt" 