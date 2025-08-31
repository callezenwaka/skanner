# Transquotation Scanner

A powerful Go-based tool for scanning and detecting quotation mark inconsistencies in your codebase. Perfect for pre-commit hooks and maintaining consistent quotation mark usage across your projects.

## Features

- 🔍 **Smart Quote Detection**: Identifies curly quotes (`"`, `"`, `'`, `'`) and suggests straight quotes
- 🌍 **International Character Support**: Recognizes legitimate uses of special characters in international text
- 📁 **Multi-file Support**: Scans entire directories with configurable include/exclude patterns
- ⚡ **Pre-commit Ready**: Designed to work seamlessly with git pre-commit hooks
- 🎯 **Configurable**: Customize scanning rules and file patterns
- 📊 **Detailed Reporting**: Provides line-by-line analysis with context

## Installation

### From Source

```bash
git clone <repository-url>
cd skanner
go build -o skanner main.go
```

### Using Go Install

```bash
go install github.com/callezenwaka/skanner@latest
```

## Usage

### Basic Usage

```bash
# Scan current directory
./skanner

# Scan specific directory
./skanner --include "**/*.go"

# Verbose output
./skanner --verbose

# Custom line length limit
./skanner --max-line-length 100
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `--include` | File patterns to include | `**/*.{go,js,ts,jsx,tsx,py,java,cpp,c,h,md,txt}` |
| `--exclude` | File patterns to exclude | `**/vendor/**,**/node_modules/**,**/.git/**` |
| `--verbose` | Enable verbose output | `false` |
| `--fix` | Automatically fix issues where possible | `false` |
| `--exit-on-error` | Exit with non-zero code on errors | `true` |
| `--max-line-length` | Maximum line length before warning | `120` |
| `--config` | Path to configuration file | (none) |

### Exit Codes

- `0`: No issues found
- `1`: Errors found (unbalanced quotes, file errors)
- `2`: Warnings found (smart quotes, line length issues)

## Pre-commit Integration

### Git Pre-commit Hook

Create a `.git/hooks/pre-commit` file:

```bash
#!/bin/sh
# Pre-commit hook for transquotation scanning

echo "🔍 Scanning for quotation mark issues..."

# Run the scanner
./skanner --exit-on-error

if [ $? -ne 0 ]; then
    echo "❌ Quotation mark issues found. Please fix them before committing."
    exit 1
fi

echo "✅ No quotation mark issues found."
exit 0
```

Make it executable:

```bash
chmod +x .git/hooks/pre-commit
```

### Pre-commit Framework

Add to your `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: local
    hooks:
      - id: skanner
        name: Transquotation Scanner
        entry: skanner
        language: system
        types: [text]
        pass_filenames: false
```

## Configuration

### Configuration File

Create a `transquotation.yaml` file:

```yaml
include_patterns:
  - "**/*.go"
  - "**/*.js"
  - "**/*.ts"

exclude_patterns:
  - "**/vendor/**"
  - "**/node_modules/**"
  - "**/.git/**"

max_line_length: 100
verbose: false
fix: false
exit_on_error: true

quotation_marks:
  - type: "smart_quotes"
    description: "Smart quotes (curly quotes)"
    pattern: "[\u201C\u201D]"
    replacement: "\""
  
  - type: "smart_single_quotes"
    description: "Smart single quotes (curly apostrophes)"
    pattern: "[\u2018\u2019]"
    replacement: "'"
```

## What It Detects

### 1. Smart Quotes
- `"` (left double quotation mark)
- `"` (right double quotation mark)
- `'` (left single quotation mark)
- `'` (right single quotation mark)

### 2. Mixed Quote Types
- Strings containing both single and double quotes
- Inconsistent quote usage

### 3. Unbalanced Quotes
- Missing opening or closing quotes
- Mismatched quote pairs

### 4. Line Length Issues
- Lines exceeding maximum length limit
- Configurable threshold

## International Character Support

The scanner recognizes legitimate uses of special characters:

- **International Text**: Allows smart quotes in strings containing non-Latin characters
- **Comments**: Permits smart quotes in documentation and comments
- **Contractions**: Allows smart apostrophes in contractions like "don't", "it's"

## Examples

### Input File
```go
package main

import "fmt"

func main() {
    // This line has "smart quotes" that should be "straight quotes"
    fmt.Println("Hello, world!")
    
    // This line is fine because it's in a comment
    fmt.Println("International text: café, naïve, résumé")
}
```

### Scanner Output
```
🔍 Found quotation mark issues in 1 files:

📁 main.go
  ⚠️ Line 6:15 - smart_quotes: Smart quotes (curly quotes)
     Context: // This line has "smart quotes" that should be "straight quotes"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the need for consistent quotation mark usage in codebases
- Built with Go's excellent text processing capabilities
- Designed for seamless integration with development workflows
