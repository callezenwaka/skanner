package main

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestInitializeQuotationMarks(t *testing.T) {
	marks := initializeQuotationMarks()

	if len(marks) == 0 {
		t.Error("Expected quotation marks to be initialized")
	}

	// Check that we have the expected types
	expectedTypes := map[string]bool{
		"smart_quotes":         false,
		"smart_single_quotes":  false,
		"backticks_in_strings": false,
		"mixed_quotes":         false,
		"unmatched_quotes":     false,
	}

	for _, mark := range marks {
		if _, exists := expectedTypes[mark.Type]; exists {
			expectedTypes[mark.Type] = true
		}
	}

	for markType, found := range expectedTypes {
		if !found {
			t.Errorf("Expected quotation mark type '%s' not found", markType)
		}
	}
}

func TestContainsInternationalText(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Hello world", false},
		{"Hello 世界", true},
		{"こんにちは", true},
		{"Привет", true},
		{"Hello", false},
	}

	for _, test := range tests {
		result := containsInternationalText(test.input)
		if result != test.expected {
			t.Errorf("containsInternationalText('%s') = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestHasBalancedQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"Hello"`, true},
		{`'Hello'`, true},
		{`"Hello'`, false},
		{`'Hello"`, false},
		{`"Hello" world`, true},
		{`"Hello world`, false},
		{`Hello world"`, false},
		{`"Hello 'world'"`, true},
		{`'Hello "world"'`, true},
		{`"Hello 'world"`, false},
		{`'Hello "world'`, false},
		{`"Hello" "world"`, true},
		{`"Hello" 'world'`, true},
		{`"Hello" 'world"`, false},
		{`Hello`, true}, // No quotes
		{`""`, true},    // Empty quotes
		{`''`, true},    // Empty single quotes
	}

	for _, test := range tests {
		result := hasBalancedQuotes(test.input)
		if result != test.expected {
			t.Errorf("hasBalancedQuotes('%s') = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsLegitimateUse(t *testing.T) {
	marks := initializeQuotationMarks()

	// Find smart quotes mark
	var smartQuotesMark QuotationMark
	for _, mark := range marks {
		if mark.Type == "smart_quotes" {
			smartQuotesMark = mark
			break
		}
	}

	if smartQuotesMark.Type == "" {
		t.Fatal("Could not find smart_quotes mark type")
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{`// This comment has "smart quotes"`, true},    // Comment
		{`/* This comment has "smart quotes" */`, true}, // Block comment
		{`# This comment has "smart quotes"`, true},     // Hash comment
		{`fmt.Println("Hello world")`, false},           // Code
		{`fmt.Println("Hello 世界")`, true},               // International text
	}

	for _, test := range tests {
		result := isLegitimateUse(test.input, smartQuotesMark)
		if result != test.expected {
			t.Errorf("isLegitimateUse('%s') = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"Hello", 10, "Hello"},
		{"Hello World", 5, "He..."},
		{"Hello World", 8, "Hello..."},
		{"", 5, ""},
		{"Hello", 0, "..."},
	}

	for _, test := range tests {
		result := truncateString(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("truncateString('%s', %d) = '%s', expected '%s'", test.input, test.maxLen, result, test.expected)
		}
	}
}

func TestFindFirstMatch(t *testing.T) {
	pattern := regexp.MustCompile(`"`)

	tests := []struct {
		input    string
		expected int
	}{
		{`Hello "world"`, 7},
		{`"Hello" world`, 1},
		{`Hello world`, 1}, // No match
		{`"Hello" "world"`, 1},
	}

	for _, test := range tests {
		result := findFirstMatch(test.input, pattern)
		if result != test.expected {
			t.Errorf("findFirstMatch('%s') = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestGetFilesToScan(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "skanner_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{
		"test.go",
		"test.js",
		"test.txt",
		"vendor/test.go",       // Should be excluded
		"node_modules/test.js", // Should be excluded
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	config := &Config{
		IncludePatterns: []string{filepath.Join(tempDir, "*.go"), filepath.Join(tempDir, "*.js")},
		ExcludePatterns: []string{filepath.Join(tempDir, "vendor"), filepath.Join(tempDir, "node_modules")},
	}

	files, err := getFilesToScan(config)
	if err != nil {
		t.Fatal(err)
	}

	// Should find test.go and test.js, but not vendor/test.go or node_modules/test.js
	expectedCount := 2
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d", expectedCount, len(files))
	}

	// Check that excluded files are not included
	for _, file := range files {
		if filepath.Base(filepath.Dir(file)) == "vendor" || filepath.Base(filepath.Dir(file)) == "node_modules" {
			t.Errorf("Excluded file found: %s", file)
		}
	}
}

func TestScanFile(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "skanner_test_*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Write test content with issues
	testContent := `package main

import "fmt"

func main() {
	// This comment has "smart quotes" which should be allowed
	fmt.Println("Hello, "world"!") // This has smart quotes that should be flagged
	fmt.Println("This is a very long line that exceeds the maximum line length limit and should trigger a warning because it's too long")
	fmt.Println("Hello world") // This is fine
}
`

	if _, err := tempFile.WriteString(testContent); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	config := &Config{
		MaxLineLength:  80,
		QuotationMarks: initializeQuotationMarks(),
	}

	result := scanFile(tempFile.Name(), config)

	if len(result.Issues) == 0 {
		t.Error("Expected to find issues in test file")
	}

	// Check that we found the expected issues
	hasSmartQuotes := false
	hasLineLength := false

	for _, issue := range result.Issues {
		switch issue.Type {
		case "smart_quotes":
			hasSmartQuotes = true
		case "line_length":
			hasLineLength = true
		}
	}

	if !hasSmartQuotes {
		t.Error("Expected to find smart quotes issue")
	}

	if !hasLineLength {
		t.Error("Expected to find line length issue")
	}
}
