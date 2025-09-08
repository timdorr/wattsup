#!/bin/bash

set -e

echo "Running tests with coverage..."

# Create coverage directory
mkdir -p coverage

# Run tests with coverage
go test -race -coverprofile=coverage/coverage.out ./...

# Generate coverage report
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Display coverage summary
go tool cover -func=coverage/coverage.out

echo "Coverage report generated at coverage/coverage.html"
