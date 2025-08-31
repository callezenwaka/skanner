package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// QuotationMark represents a type of quotation mark
type QuotationMark struct {
	Type        string
	Description string
	Pattern     *regexp.Regexp
	Replacement string
}

// ScanResult represents the result of scanning a file
type ScanResult struct {
	FilePath string
	Issues   []Issue
}

// Issue represents a quotation mark issue found in the file
type Issue struct {
	Line     int
	Column   int
	Type     string
	Message  string
	Context  string
	Severity string
}

// Configuration for the scanner
type Config struct {
	IncludePatterns []string
	ExcludePatterns []string
	QuotationMarks  []QuotationMark
	MaxLineLength   int
	Verbose         bool
	Fix             bool
	ExitOnError     bool
}

func main() {
	var (
		configPath     = flag.String("config", "", "Path to configuration file")
		includePattern = flag.String("include", "**/*.{go,js,ts,jsx,tsx,py,java,cpp,c,h,md,txt}", "File patterns to include")
		excludePattern = flag.String("exclude", "**/vendor/**,**/node_modules/**,**/.git/**", "File patterns to exclude")
		verbose        = flag.Bool("verbose", false, "Enable verbose output")
		fix            = flag.Bool("fix", false, "Automatically fix issues where possible")
		exitOnError    = flag.Bool("exit-on-error", true, "Exit with non-zero code on errors")
		maxLineLength  = flag.Int("max-line-length", 120, "Maximum line length before warning")
	)
	flag.Parse()

	config := &Config{
		IncludePatterns: strings.Split(*includePattern, ","),
		ExcludePatterns: strings.Split(*excludePattern, ","),
		MaxLineLength:   *maxLineLength,
		Verbose:         *verbose,
		Fix:             *fix,
		ExitOnError:     *exitOnError,
	}

	// Initialize quotation mark patterns
	config.QuotationMarks = initializeQuotationMarks()

	// Load custom config if provided
	if *configPath != "" {
		if err := loadConfig(*configPath, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
	}

	// Get files to scan
	files, err := getFilesToScan(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting files to scan: %v\n", err)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Printf("Found %d files to scan\n", len(files))
	}

	// Scan files
	var allResults []ScanResult
	var totalIssues int
	var hasErrors bool

	for _, file := range files {
		result := scanFile(file, config)
		if len(result.Issues) > 0 {
			allResults = append(allResults, result)
			totalIssues += len(result.Issues)

			// Check if any issues are errors
			for _, issue := range result.Issues {
				if issue.Severity == "error" {
					hasErrors = true
				}
			}
		}
	}

	// Print results
	printResults(allResults, config)

	if config.Verbose {
		fmt.Printf("\nTotal issues found: %d\n", totalIssues)
	}

	// Exit with appropriate code
	if hasErrors && config.ExitOnError {
		os.Exit(1)
	} else if totalIssues > 0 {
		os.Exit(2)
	}
}

// initializeQuotationMarks sets up the quotation mark patterns to detect
func initializeQuotationMarks() []QuotationMark {
	return []QuotationMark{
		{
			Type:        "smart_quotes",
			Description: "Smart quotes (curly quotes)",
			Pattern:     regexp.MustCompile(`[""]`),
			Replacement: `"`,
		},
		{
			Type:        "smart_single_quotes",
			Description: "Smart single quotes (curly apostrophes)",
			Pattern:     regexp.MustCompile(`['']`),
			Replacement: `'`,
		},
		{
			Type:        "backticks_in_strings",
			Description: "Backticks in string literals (potential template literal)",
			Pattern:     regexp.MustCompile(`"[^"]*` + "`" + `[^"]*"`),
			Replacement: "",
		},
		{
			Type:        "mixed_quotes",
			Description: "Mixed quote types in same string",
			Pattern:     regexp.MustCompile(`"[^"]*'[^"]*"|'[^']*"[^']*'`),
			Replacement: "",
		},
		{
			Type:        "unmatched_quotes",
			Description: "Unmatched quotes",
			Pattern:     regexp.MustCompile(`^[^"]*"[^"]*$|^[^']*'[^']*$`),
			Replacement: "",
		},
	}
}

// getFilesToScan returns a list of files to scan based on include/exclude patterns
func getFilesToScan(config *Config) ([]string, error) {
	var files []string

	for _, pattern := range config.IncludePatterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			// Check if file should be excluded
			excluded := false
			for _, excludePattern := range config.ExcludePatterns {
				if matched, _ := filepath.Match(excludePattern, match); matched {
					excluded = true
					break
				}
			}

			if !excluded {
				// Check if it's a regular file
				if info, err := os.Stat(match); err == nil && !info.IsDir() {
					files = append(files, match)
				}
			}
		}
	}

	return files, nil
}

// scanFile scans a single file for quotation mark issues
func scanFile(filePath string, config *Config) ScanResult {
	result := ScanResult{
		FilePath: filePath,
		Issues:   []Issue{},
	}

	file, err := os.Open(filePath)
	if err != nil {
		result.Issues = append(result.Issues, Issue{
			Line:     0,
			Column:   0,
			Type:     "file_error",
			Message:  fmt.Sprintf("Could not open file: %v", err),
			Severity: "error",
		})
		return result
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var content strings.Builder

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		content.WriteString(line)
		content.WriteString("\n")

		// Check line length
		if len(line) > config.MaxLineLength {
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Column:   config.MaxLineLength + 1,
				Type:     "line_length",
				Message:  fmt.Sprintf("Line exceeds maximum length (%d > %d)", len(line), config.MaxLineLength),
				Context:  truncateString(line, 50),
				Severity: "warning",
			})
		}

		// Check for quotation mark issues
		for _, qm := range config.QuotationMarks {
			if qm.Pattern.MatchString(line) {
				// Check if this is a legitimate use (e.g., international characters)
				if !isLegitimateUse(line, qm) {
					result.Issues = append(result.Issues, Issue{
						Line:     lineNum,
						Column:   findFirstMatch(line, qm.Pattern),
						Type:     qm.Type,
						Message:  qm.Description,
						Context:  truncateString(line, 80),
						Severity: "warning",
					})
				}
			}
		}

		// Check for balanced quotes
		if !hasBalancedQuotes(line) {
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Column:   1,
				Type:     "unbalanced_quotes",
				Message:  "Unbalanced quotes detected",
				Context:  truncateString(line, 80),
				Severity: "error",
			})
		}
	}

	if err := scanner.Err(); err != nil {
		result.Issues = append(result.Issues, Issue{
			Line:     0,
			Column:   0,
			Type:     "scan_error",
			Message:  fmt.Sprintf("Error scanning file: %v", err),
			Severity: "error",
		})
	}

	return result
}

// isLegitimateUse checks if the quotation mark usage is legitimate
func isLegitimateUse(line string, qm QuotationMark) bool {
	// Check for legitimate international characters or special cases
	switch qm.Type {
	case "smart_quotes":
		// Allow smart quotes in comments or documentation
		if strings.Contains(line, "//") || strings.Contains(line, "/*") || strings.Contains(line, "#") {
			return true
		}
		// Allow smart quotes in string literals that contain international text
		if containsInternationalText(line) {
			return true
		}
	case "smart_single_quotes":
		// Allow smart apostrophes in contractions or possessives
		if strings.Contains(line, "'s") || strings.Contains(line, "'t") || strings.Contains(line, "'re") {
			return true
		}
	}
	return false
}

// containsInternationalText checks if the line contains legitimate international characters
func containsInternationalText(line string) bool {
	for _, r := range line {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
			return true
		}
	}
	return false
}

// hasBalancedQuotes checks if quotes are properly balanced
func hasBalancedQuotes(line string) bool {
	var stack []rune

	for _, r := range line {
		switch r {
		case '"', '\u201C', '\u201D':
			stack = append(stack, '"')
		case '\'', '\u2018', '\u2019':
			stack = append(stack, '\'')
		case '`':
			stack = append(stack, '`')
		}
	}

	// Check if we have an even number of each quote type
	quotes := make(map[rune]int)
	for _, r := range stack {
		quotes[r]++
	}

	for _, count := range quotes {
		if count%2 != 0 {
			return false
		}
	}

	return true
}

// findFirstMatch finds the column position of the first match
func findFirstMatch(line string, pattern *regexp.Regexp) int {
	loc := pattern.FindStringIndex(line)
	if len(loc) > 0 {
		return loc[0] + 1
	}
	return 1
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return "..."
	}
	return s[:maxLen-3] + "..."
}

// printResults prints the scan results
func printResults(results []ScanResult, config *Config) {
	if len(results) == 0 {
		fmt.Println("✓ No quotation mark issues found!")
		return
	}

	fmt.Printf("🔍 Found quotation mark issues in %d files:\n\n", len(results))

	for _, result := range results {
		fmt.Printf("📁 %s\n", result.FilePath)
		for _, issue := range result.Issues {
			severityIcon := "⚠️"
			if issue.Severity == "error" {
				severityIcon = "✗"
			}

			fmt.Printf("  %s Line %d:%d - %s: %s\n",
				severityIcon, issue.Line, issue.Column, issue.Type, issue.Message)

			if issue.Context != "" {
				fmt.Printf("     Context: %s\n", issue.Context)
			}
		}
		fmt.Println()
	}
}

// loadConfig loads configuration from a file (placeholder for future implementation)
func loadConfig(configPath string, config *Config) error {
	// TODO: Implement configuration file loading
	// This could support YAML, JSON, or TOML configuration files
	return nil
}
