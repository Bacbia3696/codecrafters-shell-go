package parser

import (
	"strings"
)

// ParseLine splits a line into arguments, respecting single and double quotes.
// Text within quotes is treated as a single argument, and the quotes are removed.
// e.g., echo 'hello world' "it's fine" foo -> ["echo", "hello world", "it's fine", "foo"]
func ParseLine(line string) []string {
	args := make([]string, 0)
	var currentArg strings.Builder
	var activeQuoteChar rune = 0 // 0 means not in quotes, '\'' or '"' means in that quote type

	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		return args
	}
	lineRunes := []rune(trimmedLine)

	justClosedEmptyQuote := false

	for i := 0; i < len(lineRunes); i++ {
		char := lineRunes[i]

		// Reset flag at the start of each character processing,
		// unless it's a space immediately following the closure of an empty quote.
		// (If it's a space AND we just closed an empty quote, we want to process that state first)
		if !(char == ' ' && justClosedEmptyQuote) {
			justClosedEmptyQuote = false
		}

		if char == '\'' || char == '"' { // A quote character is encountered
			if activeQuoteChar == 0 { // Not currently in a quote, so start one
				activeQuoteChar = char
			} else if activeQuoteChar == char { // Closing the currently active quote type
				activeQuoteChar = 0 // Exited quote mode
				if currentArg.Len() == 0 {
					justClosedEmptyQuote = true
				}
			} else { // Different quote character inside an active quote (e.g. ' inside "")
				currentArg.WriteRune(char) // Treat as a literal character
			}
		} else if char == ' ' && activeQuoteChar == 0 { // Space outside of any quote
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			} else if justClosedEmptyQuote { // An empty quote pair was just closed before this space
				args = append(args, "")
				justClosedEmptyQuote = false
			}
		} else { // Regular character, or space inside quotes
			currentArg.WriteRune(char)
		}
	}

	if currentArg.Len() > 0 || justClosedEmptyQuote || activeQuoteChar != 0 {
		args = append(args, currentArg.String())
	}
	return args
}
