# Transquotation Scanner - Installation Guide

## Quick Start

### 1. Prerequisites
- Go 1.24.3 or later
- Git (for version control)

### 2. Installation Options

#### Option A: Quick Install (Recommended)
```bash
# Clone the repository
git clone <your-repo-url>
cd skanner

# Run the installation script
./install.sh
```

#### Option B: Manual Installation
```bash
# Clone the repository
git clone <your-repo-url>
cd skanner

# Build the scanner
go build -o skanner main.go

# Make it executable
chmod +x skanner

# Test the scanner
./skanner --help
```

#### Option C: Using Make
```bash
# Clone the repository
git clone <your-repo-url>
cd skanner

# Build using make
make build

# Install globally
make install

# Run tests
make test
```

### 3. Verify Installation
```bash
# Check if the scanner works
./skanner --version

# Test on a sample file
./skanner --include "*.go" --verbose
```

## Pre-commit Hook Setup

### Automatic Setup (via install.sh)
The installation script will automatically offer to set up pre-commit hooks.

### Manual Setup
```bash
# Copy the pre-commit hook
cp pre-commit .git/hooks/

# Make it executable
chmod +x .git/hooks/pre-commit
```

### Pre-commit Framework Integration
Add to your `.pre-commit-config.yaml`:
```yaml
repos:
  - repo: local
    hooks:
      - id: skanner
        name: Skanner
        entry: skanner
        language: system
        types: [text]
        pass_filenames: false
```

## Configuration

### 1. Command Line Options
```bash
./skanner --help
```

### 2. Configuration File
Create `skanner.yaml`:
```yaml
include_patterns:
  - "**/*.go"
  - "**/*.js"
  - "**/*.ts"

exclude_patterns:
  - "**/vendor/**"
  - "**/node_modules/**"

max_line_length: 100
verbose: true
exit_on_error: true
```

### 3. Environment Variables
```bash
export SKANNER_CONFIG=/path/to/config.yaml
export SKANNER_VERBOSE=true
```

## Usage Examples

### Basic Scanning
```bash
# Scan current directory
./skanner

# Scan specific files
./skanner --include "**/*.go"

# Verbose output
./skanner --verbose

# Custom line length
./skanner --max-line-length 80
```

### Integration with CI/CD
```bash
# Exit on any issues
./skanner --exit-on-error

# Scan only staged files (for pre-commit)
./skanner --include "$(git diff --cached --name-only | tr '\n' ',')"
```

### Docker Integration
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o skanner main.go

FROM alpine:latest
COPY --from=builder /app/skanner /usr/local/bin/
ENTRYPOINT ["skanner"]
```

## Troubleshooting

### Common Issues

#### 1. Permission Denied
```bash
chmod +x skanner
```

#### 2. Go Not Found
```bash
# Install Go
# Visit: https://golang.org/doc/install
```

#### 3. Build Failures
```bash
# Clean and rebuild
make clean
make build

# Check Go version
go version
```

#### 4. Pre-commit Hook Not Working
```bash
# Check if hook is executable
ls -la .git/hooks/pre-commit

# Test hook manually
.git/hooks/pre-commit

# Check git configuration
git config core.hooksPath
```

### Getting Help
```bash
# Show help
./skanner --help

# Run tests
make test

# Check version
./skanner --version
```

## Advanced Configuration

### Custom Quotation Mark Patterns
```yaml
quotation_marks:
  - type: "custom_quotes"
    description: "Custom quote pattern"
    pattern: "[\\u201C\\u201D]"
    replacement: "\""
```

### Language-Specific Rules
```yaml
language_rules:
  go:
    allow_smart_quotes_in_comments: true
    allow_smart_quotes_in_strings: false
  
  markdown:
    allow_smart_quotes: true
    allow_smart_apostrophes: true
```

### Performance Tuning
```yaml
# For large codebases
max_concurrent_scans: 4
buffer_size: 8192
skip_binary_files: true
```

## Contributing

### Development Setup
```bash
# Clone and setup
git clone <repo-url>
cd skanner

# Install development dependencies
go mod download

# Run tests
make test

# Run linting
make lint

# Run all checks
make check
```

### Running Tests
```bash
# Run all tests
go test -v

# Run specific test
go test -v -run TestScanFile

# Run with coverage
make test-coverage
```

## Support

- **Documentation**: README.md
- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions
- **Wiki**: GitHub Wiki

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
