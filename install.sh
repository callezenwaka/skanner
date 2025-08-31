#!/bin/bash

# Transquotation Scanner Installation Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Transquotation Scanner - Installation${NC}"

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}✗ Go is not installed. Please install Go first.${NC}"
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed: $(go version)${NC}"

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Build the scanner
echo -e "${BLUE}🔨 Building skanner...${NC}"
go build -o skanner main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Make the scanner executable
chmod +x skanner

# Ask user if they want to install globally
read -p "Do you want to install the scanner globally? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    # Install to /usr/local/bin
    if [ -w /usr/local/bin ]; then
        sudo cp skanner /usr/local/bin/
        echo -e "${GREEN}✓ Scanner installed to /usr/local/bin/skanner${NC}"
    else
        echo -e "${YELLOW}⚠️  Cannot write to /usr/local/bin. Installing to ~/.local/bin instead.${NC}"
        mkdir -p ~/.local/bin
        cp skanner ~/.local/bin/
        echo -e "${GREEN}✓ Scanner installed to ~/.local/bin/skanner${NC}"
        echo -e "${YELLOW}⚠️  Please add ~/.local/bin to your PATH if not already added.${NC}"
    fi
else
    echo -e "${BLUE}📁 Scanner is available at: $SCRIPT_DIR/skanner${NC}"
fi

# Ask user if they want to set up pre-commit hooks
read -p "Do you want to set up pre-commit hooks for this repository? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    # Check if this is a git repository
    if [ -d .git ]; then
        echo -e "${BLUE}🔧 Setting up pre-commit hook...${NC}"
        
        # Copy the pre-commit hook
        cp pre-commit .git/hooks/
        chmod +x .git/hooks/pre-commit

        echo -e "${GREEN}✓ Pre-commit hook installed${NC}"
        echo -e "${BLUE}📝 The hook will now run automatically before each commit${NC}"
    else
        echo -e "${YELLOW}⚠️  This directory is not a git repository. Skipping pre-commit hook setup.${NC}"
    fi
fi

# Test the scanner
echo -e "${BLUE}🧪 Testing the scanner...${NC}"
./skanner --help

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Scanner is working correctly${NC}"
else
    echo -e "${RED}✗ Scanner test failed${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Installation complete!${NC}"
echo ""
echo -e "${BLUE}Usage examples:${NC}"
echo "  ./skanner                    # Scan current directory"
echo "  ./skanner --verbose          # Verbose output"
echo "  ./skanner --help             # Show help"
echo ""
echo -e "${BLUE}For more information, see:${NC}"
echo "  README.md"
echo "  ./skanner --help"
echo ""
echo -e "${BLUE}To run tests:${NC}"
echo "  make test"
echo ""
echo -e "${BLUE}To build for other platforms:${NC}"
echo "  make build-all"
