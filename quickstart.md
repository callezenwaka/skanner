# Navigate to the project
```bash
cd skanner

# Build the scanner
make build

# Test it
./build/skanner --help

# Scan your code
./build/skanner --verbose

# Set up pre-commit hooks
./install.sh
```

# Build and Copy the Binary
```bash
# Build the scanner
cd /Users/callisezenwaka/Documents/projects/demos/skanner
make build

# Copy to your other project
cp build/skanner /path/to/your/other/project/

# Use it in that project
cd /path/to/your/other/project
./skanner --verbose
```

# Create a Script in Each Project
```bash
# In your other project, create a scan script
cat > scan-quotes.sh << 'EOF'
#!/bin/bash
# Script to run transquotation scanner

if [ -f "./skanner" ]; then
    ./skanner "$@"
else
    echo "Transquotation scanner not found. Please copy it to this project first."
    exit 1
fi
EOF

chmod +x scan-quotes.sh

# Use it
./scan-quotes.sh --verbose
```

# Scan a Go Project
```bash
# Navigate to your Go project
cd /path/to/my-go-project

# Use the globally installed scanner
transquotation-scanner --include "**/*.go" --verbose

# Or use go run
go run /Users/callisezenwaka/Documents/projects/demos/golang/transquotation-scanner/main.go --include "**/*.go" --verbose
```

# Scan a JavaScript Project
```bash
# Navigate to your JS project
cd /path/to/my-js-project

# Scan JavaScript files
transquotation-scanner --include "**/*.{js,jsx,ts,tsx}" --verbose
```

# Scan a Mixed Language Project
```bash
# Navigate to your project
cd /path/to/my-project

# Scan multiple file types
transquotation-scanner \
  --include "**/*.{go,js,ts,py,java,md}" \
  --exclude "**/node_modules/**,**/vendor/**" \
  --verbose
```

# Set Up Pre-commit Hooks
```bash
# In your other project
cd /path/to/your/other/project

# Copy the pre-commit hook
cp /Users/callisezenwaka/Documents/projects/demos/golang/transquotation-scanner/pre-commit .git/hooks/

# Make it executable
chmod +x .git/hooks/pre-commit

# Edit the hook to point to your scanner location
# Change this line in .git/hooks/pre-commit:
# SCANNER_PATH="$HOOK_DIR/transquotation-scanner"
# To:
# SCANNER_PATH="/usr/local/bin/transquotation-scanner"
# Or wherever you installed it
```

# Running Tests
## Run All Tests
```bash
# Basic test run
go test

# Verbose output (shows each test)
go test -v

# Run with coverage
go test -cover
```

## Run Specific Tests
```bash
# Run only one test function
go test -v -run TestScanFile

# Run tests matching a pattern
go test -v -run Test.*Quotes

# Run tests in a specific file
go test -v scanner_test.go
```

## Test Coverage
```bash
go test -cover
go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
```