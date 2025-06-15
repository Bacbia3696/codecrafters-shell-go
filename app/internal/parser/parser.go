package parser

import (
	"strings"

	shellerrors "github.com/codecrafters-io/shell-starter-go/app/internal/errors"
)

// RedirectionType defines the type of output redirection.
// NoRedirection indicates no redirection.
// StdoutRedirection indicates standard output redirection ('>' or '1>').
// StderrRedirection indicates standard error redirection ('2>').
// StdoutAppendRedirection indicates standard output append redirection ('>>').
// StderrAppendRedirection indicates standard error append redirection ('2>>').
const (
	NoRedirection = iota
	StdoutRedirection
	StderrRedirection
	StdoutAppendRedirection
	StderrAppendRedirection
)

// RedirectionInfo holds information about redirection found in a command
type RedirectionInfo struct {
	Type     int
	Filename string
	Found    bool
}

// QuoteTracker helps track quote state while parsing
type QuoteTracker struct {
	activeQuoteChar rune
}

// updateQuoteState updates the quote tracking state
func (qt *QuoteTracker) updateQuoteState(char rune) {
	if char == '\'' || char == '"' {
		if qt.activeQuoteChar == 0 {
			qt.activeQuoteChar = char
		} else if qt.activeQuoteChar == char {
			qt.activeQuoteChar = 0
		}
	}
}

// isInsideQuotes returns true if currently inside quotes
func (qt *QuoteTracker) isInsideQuotes() bool {
	return qt.activeQuoteChar != 0
}

// isStandaloneDigit checks if a digit at position i is standalone (preceded by space or at start)
func isStandaloneDigit(runes []rune, i int) bool {
	return i == 0 || runes[i-1] == ' '
}

// isGluedToArgument checks if a redirection operator is glued to an argument (no space before)
func isGluedToArgument(runes []rune, i int) bool {
	if i-1 >= 0 && runes[i-1] >= '0' && runes[i-1] <= '9' {
		return i-2 < 0 || runes[i-2] != ' '
	}
	return false
}

// isGluedToArgumentGeneral checks if a redirection operator is glued to any non-space character
func isGluedToArgumentGeneral(runes []rune, i int) bool {
	return i-1 >= 0 && runes[i-1] != ' '
}

// createRedirectionInfo creates a RedirectionInfo struct from parsed components
func createRedirectionInfo(redirectType int, filename string) RedirectionInfo {
	return RedirectionInfo{
		Type:     redirectType,
		Filename: stripQuotes(filename),
		Found:    true,
	}
}

// tryParseNumberedRedirection attempts to parse '1>' or '2>' redirection
func tryParseNumberedRedirection(runes []rune, i int, expectedDigit rune, redirectType int) (commandPart string, redirect RedirectionInfo, found bool) {
	n := len(runes)

	// Check for pattern: digit + '>'
	if i+1 < n && runes[i] == expectedDigit && runes[i+1] == '>' {
		if isStandaloneDigit(runes, i) {
			commandPart = strings.TrimSpace(string(runes[:i]))
			filename := strings.TrimSpace(string(runes[i+2:]))
			return commandPart, createRedirectionInfo(redirectType, filename), true
		}
	}
	return "", RedirectionInfo{}, false
}

// tryParseNumberedAppendRedirection attempts to parse '1>>' or '2>>' redirection
func tryParseNumberedAppendRedirection(runes []rune, i int, expectedDigit rune, redirectType int) (commandPart string, redirect RedirectionInfo, found bool) {
	n := len(runes)

	// Check for pattern: digit + '>>'
	if i+2 < n && runes[i] == expectedDigit && runes[i+1] == '>' && runes[i+2] == '>' {
		if isStandaloneDigit(runes, i) {
			commandPart = strings.TrimSpace(string(runes[:i]))
			filename := strings.TrimSpace(string(runes[i+3:]))
			return commandPart, createRedirectionInfo(redirectType, filename), true
		}
	}
	return "", RedirectionInfo{}, false
}

// parseGenericRedirection handles generic '>' redirection with optional file descriptor
func parseGenericRedirection(runes []rune, i int) (commandPart string, redirect RedirectionInfo) {
	// Skip if redirection is glued to an argument
	if isGluedToArgument(runes, i) {
		return "", RedirectionInfo{}
	}

	// Walk backwards beyond spaces to find potential file descriptor
	j := i - 1
	for j >= 0 && runes[j] == ' ' {
		j--
	}

	redirectType := StdoutRedirection
	commandEndIndex := j + 1

	// Check for explicit file descriptor (1 or 2)
	if j >= 0 && (runes[j] == '1' || runes[j] == '2') {
		if isStandaloneDigit(runes, j) {
			if runes[j] == '2' {
				redirectType = StderrRedirection
			}
			commandEndIndex = j
		}
	}

	commandPart = strings.TrimSpace(string(runes[:commandEndIndex]))
	filename := strings.TrimSpace(string(runes[i+1:]))
	return commandPart, createRedirectionInfo(redirectType, filename)
}

// parseGenericAppendRedirection handles generic '>>' redirection with optional file descriptor
func parseGenericAppendRedirection(runes []rune, i int) (commandPart string, redirect RedirectionInfo) {
	n := len(runes)

	// Make sure we have '>>' pattern
	if i+1 >= n || runes[i+1] != '>' {
		return "", RedirectionInfo{}
	}

	// Skip if redirection is glued to an argument (for >> we check any non-space character)
	if isGluedToArgumentGeneral(runes, i) {
		return "", RedirectionInfo{}
	}

	// Walk backwards beyond spaces to find potential file descriptor
	j := i - 1
	for j >= 0 && runes[j] == ' ' {
		j--
	}

	redirectType := StdoutAppendRedirection
	commandEndIndex := j + 1

	// Check for explicit file descriptor (1 or 2)
	if j >= 0 && (runes[j] == '1' || runes[j] == '2') {
		if isStandaloneDigit(runes, j) {
			if runes[j] == '2' {
				redirectType = StderrAppendRedirection
			}
			commandEndIndex = j
		}
	}

	commandPart = strings.TrimSpace(string(runes[:commandEndIndex]))
	filename := strings.TrimSpace(string(runes[i+2:]))
	return commandPart, createRedirectionInfo(redirectType, filename)
}

// findRedirection scans the line for the first unquoted redirection operator
// and returns information about the redirection found
func findRedirection(line string) (commandPart string, redirect RedirectionInfo) {
	runes := []rune(line)
	n := len(runes)
	quoteTracker := &QuoteTracker{}

	for i := range n {
		char := runes[i]

		// Update quote state
		quoteTracker.updateQuoteState(char)

		// Only look for redirection operators if not inside quotes
		if !quoteTracker.isInsideQuotes() {
			// Try parsing '2>>' redirection first (longer pattern)
			if cmdPart, redir, found := tryParseNumberedAppendRedirection(runes, i, '2', StderrAppendRedirection); found {
				return cmdPart, redir
			}

			// Try parsing '1>>' redirection
			if cmdPart, redir, found := tryParseNumberedAppendRedirection(runes, i, '1', StdoutAppendRedirection); found {
				return cmdPart, redir
			}

			// Try parsing '2>' redirection
			if cmdPart, redir, found := tryParseNumberedRedirection(runes, i, '2', StderrRedirection); found {
				return cmdPart, redir
			}

			// Try parsing '1>' redirection
			if cmdPart, redir, found := tryParseNumberedRedirection(runes, i, '1', StdoutRedirection); found {
				return cmdPart, redir
			}

			// Try parsing generic '>>' or '>' redirection
			if char == '>' {
				// Check for '>>' first (longer pattern)
				if i+1 < n && runes[i+1] == '>' {
					// Try to parse '>>' - if it succeeds, return it
					if cmdPart, redir := parseGenericAppendRedirection(runes, i); redir.Found {
						return cmdPart, redir
					}
					// If '>>' was detected but rejected (e.g., glued to argument), skip both '>' characters
					// This prevents "hello>>out.txt" from being parsed as "hello>" with redirection
					// We need to skip the next iteration too since we've already processed both '>' characters
					i++ // Skip the second '>' character
					continue
				}

				// Single '>' redirection (only if not part of '>>')
				if cmdPart, redir := parseGenericRedirection(runes, i); redir.Found {
					return cmdPart, redir
				}
			}
		}
	}

	return line, RedirectionInfo{Found: false}
}

// tokenize splits a command string into arguments, respecting quotes
func tokenize(commandStr string) ([]string, error) {
	if strings.TrimSpace(commandStr) == "" {
		return []string{}, nil
	}

	var args []string
	var currentArg strings.Builder
	var activeQuoteChar rune = 0
	justClosedEmptyQuote := false

	runes := []rune(strings.TrimSpace(commandStr))

	for i := range runes {
		char := runes[i]

		// Reset empty quote flag unless we're handling a space after empty quote
		if !(char == ' ' && justClosedEmptyQuote) {
			justClosedEmptyQuote = false
		}

		// Handle quote characters
		if char == '\'' || char == '"' {
			if activeQuoteChar == 0 {
				// Start quote
				activeQuoteChar = char
			} else if activeQuoteChar == char {
				// End quote
				activeQuoteChar = 0
				if currentArg.Len() == 0 {
					justClosedEmptyQuote = true
				}
			} else {
				// Different quote inside active quote - treat as literal
				currentArg.WriteRune(char)
			}
		} else if char == ' ' && activeQuoteChar == 0 {
			// Space outside quotes - end current argument
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			} else if justClosedEmptyQuote {
				args = append(args, "")
				justClosedEmptyQuote = false
			}
		} else {
			// Regular character or space inside quotes
			currentArg.WriteRune(char)
		}
	}

	// Add final argument if any
	if currentArg.Len() > 0 || justClosedEmptyQuote || activeQuoteChar != 0 {
		args = append(args, currentArg.String())
	}

	return args, nil
}

// validateRedirection checks if redirection parameters are valid
func validateRedirection(redirect RedirectionInfo) error {
	if redirect.Found && redirect.Filename == "" {
		return shellerrors.NewParseError("missing filename for redirection")
	}
	return nil
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

// ParseLine splits a line into arguments and redirection information if present.
// It handles '>', '>>', '1>', '1>>', '2>', and '2>>' output redirection operators.
// Text within quotes is treated as a single argument, and the quotes are removed.
// e.g., echo 'hello world' > out.txt -> args=["echo", "hello world"], outputFile="out.txt", append=false
// e.g., echo 'hello world' >> out.txt -> args=["echo", "hello world"], outputFile="out.txt", append=true
// e.g., ls /foo 2> err.txt -> args=["ls", "/foo"], errorFile="err.txt", append=false
func ParseLine(line string) (args []string, outputFile string, errorFile string, err error) {
	args, outputFile, errorFile, _, _, err = ParseLineWithMode(line)
	return args, outputFile, errorFile, err
}

// ParseLineWithMode splits a line into arguments and redirection information with append mode.
// It handles '>', '>>', '1>', '1>>', '2>', and '2>>' output redirection operators.
// Text within quotes is treated as a single argument, and the quotes are removed.
// Returns append mode flags for both stdout and stderr redirection.
func ParseLineWithMode(line string) (args []string, outputFile string, errorFile string, outputAppend bool, errorAppend bool, err error) {
	// Handle empty input
	if strings.TrimSpace(line) == "" {
		return []string{}, "", "", false, false, nil
	}

	// Find and extract redirection
	commandPart, redirect := findRedirection(line)

	// Validate redirection
	if err := validateRedirection(redirect); err != nil {
		return nil, "", "", false, false, err
	}

	// Set output files and append modes based on redirection type
	if redirect.Found {
		switch redirect.Type {
		case StdoutRedirection:
			outputFile = redirect.Filename
			outputAppend = false
		case StdoutAppendRedirection:
			outputFile = redirect.Filename
			outputAppend = true
		case StderrRedirection:
			errorFile = redirect.Filename
			errorAppend = false
		case StderrAppendRedirection:
			errorFile = redirect.Filename
			errorAppend = true
		}
	}

	// Tokenize the command part
	args, err = tokenize(commandPart)
	if err != nil {
		return nil, "", "", false, false, err
	}

	return args, outputFile, errorFile, outputAppend, errorAppend, nil
}
