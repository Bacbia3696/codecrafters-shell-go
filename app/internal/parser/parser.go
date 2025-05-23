package parser

import (
	"strings"

	shellerrors "github.com/codecrafters-io/shell-starter-go/app/internal/errors"
)

// RedirectionType defines the type of output redirection.
// NoRedirection indicates no redirection.
// StdoutRedirection indicates standard output redirection ('>' or '1>').
// StderrRedirection indicates standard error redirection ('2>').
const (
	NoRedirection = iota
	StdoutRedirection
	StderrRedirection
)

// splitByRedirect scans the line for the first unquoted redirection operator ('>', '1>', or '2>')
// and splits the line into a command part, a filename part, and the type of redirection.
// It returns the command part, filename part, redirection type, and a boolean indicating if redirection was found.
func splitByRedirect(line string) (commandPart string, filenamePart string, redirectType int, foundRedirect bool) {
	runes := []rune(line)
	n := len(runes)
	activeQuoteChar := rune(0)

	// Scan for '1>' or '2>' first, then for '>'.
	// This prioritizes specific stderr/stdout redirection over general output redirection.
	for i := 0; i < n; i++ {
		char := runes[i]
		if char == '\'' || char == '"' { // Handle quotes
			if activeQuoteChar == 0 {
				activeQuoteChar = char
			} else if activeQuoteChar == char {
				activeQuoteChar = 0
			}
			continue
		}

		if activeQuoteChar == 0 { // Only look for redirection operators if not inside quotes
			// Check for '2>'
			if char == '2' && i+1 < n && runes[i+1] == '>' {
				if i == 0 || runes[i-1] == ' ' { // Ensure '2' is standalone or preceded by space
					commandPart = strings.TrimSpace(string(runes[:i]))
					filenamePart = strings.TrimSpace(string(runes[i+2:]))
					// Remove quotes from filename if present
					filenamePart = stripQuotes(filenamePart)
					return commandPart, filenamePart, StderrRedirection, true
				}
			}
			// Check for '1>'
			if char == '1' && i+1 < n && runes[i+1] == '>' {
				if i == 0 || runes[i-1] == ' ' { // Ensure '1' is standalone or preceded by space
					commandPart = strings.TrimSpace(string(runes[:i]))
					filenamePart = strings.TrimSpace(string(runes[i+2:]))
					filenamePart = stripQuotes(filenamePart)
					return commandPart, filenamePart, StdoutRedirection, true
				}
			}
			// Check for '>'
			if char == '>' {
				commandPart = strings.TrimSpace(string(runes[:i]))
				filenamePart = strings.TrimSpace(string(runes[i+1:]))
				filenamePart = stripQuotes(filenamePart)
				return commandPart, filenamePart, StdoutRedirection, true // Default '>' is for stdout
			}
		}
	}

	return line, "", NoRedirection, false // No redirect operator found, commandPart is the original line
}

// stripQuotes removes a single layer of leading and trailing quotes (' or ") if present.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// ParseLine splits a line into arguments and an output filename if redirection is present.
// It handles '>', '1>' and '2>' output redirection operators.
// Text within quotes is treated as a single argument, and the quotes are removed.
// e.g., echo 'hello world' > out.txt -> args=["echo", "hello world"], outputFile="out.txt"
// e.g., ls /foo 2> err.txt -> args=["ls", "/foo"], errorFile="err.txt"
func ParseLine(line string) (args []string, outputFile string, errorFile string, err error) {
	args = make([]string, 0) // Ensure args is initialized

	trimmedOriginalLine := strings.TrimSpace(line)
	if trimmedOriginalLine == "" {
		return args, "", "", nil // Empty line results in no arguments and no redirection
	}

	commandPartStr, filenameStr, redirectType, redirectFound := splitByRedirect(trimmedOriginalLine)

	if redirectFound {
		if filenameStr == "" {
			return nil, "", "", shellerrors.NewParseError("missing filename for redirection")
		}
		switch redirectType {
		case StdoutRedirection:
			outputFile = filenameStr
		case StderrRedirection:
			errorFile = filenameStr
		}
	}

	// If no redirect, commandPartStr is trimmedOriginalLine, and outputFile/errorFile are ""

	// Proceed to parse commandPartStr using the existing argument parsing logic
	var currentArg strings.Builder
	var activeQuoteChar rune = 0 // 0 means not in quotes, '\'' or '"' means in that quote type
	justClosedEmptyQuote := false

	// If commandPartStr is empty (e.g., line was "> out.txt"), no args to parse.
	if strings.TrimSpace(commandPartStr) == "" {
		return args, outputFile, errorFile, nil
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
	return args, outputFile, errorFile, nil
}
