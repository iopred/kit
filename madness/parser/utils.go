package main

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidateInput checks if the input has proper bracket matching and basic syntax
func ValidateInput(input string) error {
	bracketCount := 0
	lines := strings.Split(input, "\n")

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		// Count brackets
		bracketCount += strings.Count(line, "{")
		bracketCount -= strings.Count(line, "}")

		// Check for common syntax errors
		if strings.HasPrefix(line, "}") && strings.Contains(line, "{") {
			return fmt.Errorf("line %d: invalid bracket placement - closing bracket cannot be followed by opening bracket", lineNum+1)
		}

		// Validate conversation format
		if strings.Contains(line, ">:") {
			if !strings.HasPrefix(line, "<") {
				return fmt.Errorf("line %d: invalid conversation format - must start with '<'", lineNum+1)
			}
		}
	}

	if bracketCount != 0 {
		return fmt.Errorf("unmatched brackets: %d unclosed brackets", bracketCount)
	}

	return nil
}

// DebugPrint prints the current state of parsing for debugging
func DebugPrint(p *Parser) {
	fmt.Printf("Current parsing state:\n")
	fmt.Printf("  Current Node: %s\n", p.currentNode)
	fmt.Printf("  Bracket Depth: %d\n", p.bracketDepth)
	fmt.Printf("  Risk Cost: %f\n", p.riskCost)
}

// IsTimeCoordinate checks if a string represents a time coordinate (e.g., "0000", "1000")
func IsTimeCoordinate(s string) bool {
	if len(s) != 4 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// SanitizeNodeName removes any invalid characters from node names
func SanitizeNodeName(name string) string {
	// Allow emojis and alphanumeric characters
	var result strings.Builder
	for _, r := range name {
		if r >= 0x1F300 || // Emoji range
			unicode.IsLetter(r) ||
			unicode.IsNumber(r) ||
			r == ' ' || r == '_' {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
} 