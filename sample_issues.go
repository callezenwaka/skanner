package main

import (
	"fmt"
)

// This file demonstrates various quotation mark issues
// that the transquotation scanner can detect

func demonstrateIssues() {
	// ✗ Smart quotes (curly quotes) - should be straight quotes
	fmt.Println("Hello, \"world\"!")

	// ✗ Smart single quotes (curly apostrophes) - should be straight apostrophes
	fmt.Println("Don't worry, it's fine")

	// ✗ Mixed quote types in same string
	fmt.Println("This string has \"quotes\" and 'quotes' mixed")

	// ✗ Unmatched quotes - this will cause a syntax error in Go
	// fmt.Println("This string is missing a closing quote

	// ✓ This line is fine - straight quotes
	fmt.Println("This is correct")

	// ✓ This line is fine - smart quotes in comment
	// This comment has "smart quotes" which are allowed

	// ✓ This line is fine - international text with smart quotes
	fmt.Println("International text: café, naïve, résumé")

	// ✗ Line length issue (if max-line-length is set to 80)
	fmt.Println("This is a very long line that exceeds the maximum line length limit and should trigger a warning")

	// ✗ Backticks in string literals
	fmt.Println("This string has `backticks` inside")

	// ✓ Template literals are fine in some contexts
	message := `This is a template literal
	with multiple lines
	and it's perfectly fine`

	fmt.Println(message)

	// ✗ Unbalanced quotes across multiple lines - this will cause syntax error
	// unbalanced := "This string starts here
	// but doesn't end properly

	// ✓ Properly balanced quotes
	balanced := "This string is properly balanced"

	// ✗ Smart quotes in string literals (not in comments)
	smartQuotes := "This has \"smart quotes\" that should be flagged"

	// ✓ Smart quotes in comments are fine
	// This comment: "smart quotes" are allowed

	// ✗ Mixed quote types in variable names (if scanning for this)
	// var `backtick` = "value"  // This would cause syntax error

	fmt.Println(balanced)
	fmt.Println(smartQuotes)
	// fmt.Println(`backtick`)
}
