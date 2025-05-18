package parser

import (
	"errors"
	"strings"
)

// splitByRedirect scans the line for the first unquoted redirection operator ('>' or '1>')
// and splits the line into a command part and a filename part.
// It returns the command part, filename part, and a boolean indicating if redirection was found.
func splitByRedirect(line string) (commandPart string, filenamePart string, foundRedirect bool) {
	runes := []rune(line)
	n := len(runes)
	activeQuoteChar := rune(0)

	// First, try to find "1>" where '1' is standalone
	for i := 0; i < n-1; i++ { // Iterate up to n-2 to check runes[i] and runes[i+1]
		char := runes[i]
		if char == '\'' || char == '"' {
			if activeQuoteChar == 0 {
				activeQuoteChar = char
			} else if activeQuoteChar == char {
				activeQuoteChar = 0
			}
		} else if activeQuoteChar == 0 && char == '1' && runes[i+1] == '>' {
			// Found "1>". Check if '1' is standalone (start of line or preceded by space).
			if i == 0 || runes[i-1] == ' ' {
				commandPart = strings.TrimSpace(string(runes[:i]))
				filenamePart = strings.TrimSpace(string(runes[i+2:]))
				return commandPart, filenamePart, true
			}
		}
	}

	// If "1>" wasn't found or wasn't applicable, try to find ">"
	activeQuoteChar = rune(0) // Reset quote tracking for this scan
	for i := 0; i < n; i++ {
		char := runes[i]
		if char == '\'' || char == '"' {
			if activeQuoteChar == 0 {
				activeQuoteChar = char
			} else if activeQuoteChar == char {
				activeQuoteChar = 0
			}
		} else if activeQuoteChar == 0 && char == '>' {
			// Found ">"
			commandPart = strings.TrimSpace(string(runes[:i]))
			filenamePart = strings.TrimSpace(string(runes[i+1:]))
			return commandPart, filenamePart, true
		}
	}

	return line, "", false // No redirect operator found, commandPart is the original line
}

// ParseLine splits a line into arguments and an output filename if redirection is present.
// It handles '>' and '1>' output redirection operators.
// Text within quotes is treated as a single argument, and the quotes are removed.
// e.g., echo 'hello world' > out.txt -> args=["echo", "hello world"], outputFile="out.txt"
func ParseLine(line string) (args []string, outputFile string, err error) {
	args = make([]string, 0) // Ensure args is initialized

	trimmedOriginalLine := strings.TrimSpace(line)
	if trimmedOriginalLine == "" {
		return args, "", nil // Empty line results in no arguments and no redirection
	}

	commandPartStr, filenameStr, redirectFound := splitByRedirect(trimmedOriginalLine)

	if redirectFound {
		if filenameStr == "" {
			return nil, "", errors.New("missing filename for redirection")
		}
		outputFile = filenameStr
		// commandPartStr now contains the command and its arguments to be parsed
	}
	// If no redirect, commandPartStr is trimmedOriginalLine, and outputFile is ""

	// Proceed to parse commandPartStr using the existing argument parsing logic
	var currentArg strings.Builder
	var activeQuoteChar rune = 0 // 0 means not in quotes, '\'' or '"' means in that quote type
	justClosedEmptyQuote := false

	// If commandPartStr is empty (e.g., line was "> out.txt"), no args to parse.
	if strings.TrimSpace(commandPartStr) == "" {
		return args, outputFile, nil
	}

	lineRunes := []rune(strings.TrimSpace(commandPartStr)) // Parse the command part

	for i := range len(lineRunes) {
		char := lineRunes[i]

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

	// Add the last argument if any, or if an unclosed quote exists (maintaining original behavior)
	if currentArg.Len() > 0 || justClosedEmptyQuote || activeQuoteChar != 0 {
		args = append(args, currentArg.String())
	}
	return args, outputFile, nil
}
