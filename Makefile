# Makefile for Transquotation Scanner

# Variables
BINARY_NAME=skanner
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "🔨 Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} main.go
	@echo "✓ Build complete: ${BUILD_DIR}/${BINARY_NAME}"

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "🔨 Building for Linux..."
	@mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 main.go
	@echo "✓ Linux build complete: ${BUILD_DIR}/${BINARY_NAME}-linux-amd64"

# Build for macOS
.PHONY: build-darwin
build-darwin:
	@echo "🔨 Building for macOS..."
	@mkdir -p ${BUILD_DIR}
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 main.go
	@echo "✓ macOS builds complete"

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "🔨 Building for Windows..."
	@mkdir -p ${BUILD_DIR}
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe main.go
	@echo "✓ Windows build complete: ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe"

# Install the binary
.PHONY: install
install: build
	@echo "📦 Installing ${BINARY_NAME}..."
	cp ${BUILD_DIR}/${BINARY_NAME} /usr/local/bin/
	@echo "✓ Installation complete"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf ${BUILD_DIR}
	@echo "✓ Clean complete"

# Run tests
.PHONY: test
test:
	@echo "🧪 Running tests..."
	go test -v ./...
	@echo "✓ Tests complete"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Run the scanner on the current directory
.PHONY: scan
scan: build
	@echo "🔍 Running scanner on current directory..."
	./${BUILD_DIR}/${BINARY_NAME} --verbose

# Run the scanner with custom options
.PHONY: scan-custom
scan-custom: build
	@echo "🔍 Running scanner with custom options..."
	./${BUILD_DIR}/${BINARY_NAME} --include "**/*.go" --max-line-length 100 --verbose

# Format code
.PHONY: fmt
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✓ Formatting complete"

# Lint code
.PHONY: lint
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️ golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Run all checks
.PHONY: check
check: fmt lint test
	@echo "✓ All checks passed"

# Create release archive
.PHONY: release
release: build-all
	@echo "📦 Creating release archive..."
	@mkdir -p ${BUILD_DIR}/release
	cd ${BUILD_DIR} && tar -czf release/${BINARY_NAME}-${VERSION}.tar.gz \
		${BINARY_NAME}-linux-amd64 \
		${BINARY_NAME}-darwin-amd64 \
		${BINARY_NAME}-darwin-arm64 \
		${BINARY_NAME}-windows-amd64.exe
	@echo "✓ Release archive created: ${BUILD_DIR}/release/${BINARY_NAME}-${VERSION}.tar.gz"

# Show help
.PHONY: help
help:
	@echo "Transquotation Scanner - Available targets:"
	@echo ""
	@echo "  build          - Build the binary for current platform"
	@echo "  build-all      - Build for Linux, macOS, and Windows"
	@echo "  install        - Install the binary to /usr/local/bin"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  scan           - Run scanner on current directory"
	@echo "  scan-custom    - Run scanner with custom options"
	@echo "  fmt            - Format Go code"
	@echo "  lint           - Lint Go code"
	@echo "  check          - Run all checks (fmt, lint, test)"
	@echo "  release        - Create release archive"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make scan-custom"
	@echo "  make build-all"
	@echo "  make install"
