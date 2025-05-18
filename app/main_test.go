package main

import (
	"reflect"
	"testing"
)

func TestParseLineWithQuotes(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		expectedArgs []string
	}{
		{
			name:         "empty string",
			line:         "",
			expectedArgs: []string{},
		},
		{
			name:         "only spaces",
			line:         "   ",
			expectedArgs: []string{},
		},
		{
			name:         "simple command",
			line:         "echo",
			expectedArgs: []string{"echo"},
		},
		{
			name:         "command with one arg",
			line:         "echo hello",
			expectedArgs: []string{"echo", "hello"},
		},
		{
			name:         "command with multiple args",
			line:         "ls -l /tmp",
			expectedArgs: []string{"ls", "-l", "/tmp"},
		},
		{
			name:         "command with single quoted arg",
			line:         "echo 'hello world'",
			expectedArgs: []string{"echo", "hello world"},
		},
		{
			name:         "command with empty single quoted arg",
			line:         "echo ''",
			expectedArgs: []string{"echo", ""},
		},
		{
			name:         "command with multiple single quoted args",
			line:         "echo 'hello world' 'another one'",
			expectedArgs: []string{"echo", "hello world", "another one"},
		},
		{
			name:         "mixed quoted and unquoted",
			line:         "command 'quoted arg' unquoted 'another quoted' last",
			expectedArgs: []string{"command", "quoted arg", "unquoted", "another quoted", "last"},
		},
		{
			name:         "leading and trailing spaces",
			line:         "  echo hello  ",
			expectedArgs: []string{"echo", "hello"},
		},
		{
			name:         "multiple spaces between args",
			line:         "echo   hello   world",
			expectedArgs: []string{"echo", "hello", "world"},
		},
		{
			name:         "quoted arg with leading/trailing space inside quotes",
			line:         "echo '  spaced arg  '",
			expectedArgs: []string{"echo", "  spaced arg  "},
		},
		{
			name: "adjacent quoted args",
			line: "echo 'hello''world'",
			// This behavior might be specific. Current implementation would likely produce ["echo", "helloworld"]
			// or ["echo", "hello", "world"] if the logic correctly separates based on quote boundaries
			// Let's assume the current code merges them if there's no space, or if the quote logic is as implemented.
			// Based on current parseLineWithQuotes, it should be "helloworld"
			expectedArgs: []string{"echo", "helloworld"},
		},
		{
			name:         "complex case with multiple quotes and spaces",
			line:         " cmd 'arg1 part1' '' arg2 'arg3 part1 part2' ",
			expectedArgs: []string{"cmd", "arg1 part1", "", "arg2", "arg3 part1 part2"},
		},
		{
			name:         "line ending with an open quote (should be treated as unclosed)",
			line:         "echo 'hello",
			expectedArgs: []string{"echo", "hello"},
		},
		{
			name:         "line with quote in middle of word",
			line:         "echo hel'lo",
			expectedArgs: []string{"echo", "hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualArgs := parseLineWithQuotes(tt.line)
			if !reflect.DeepEqual(actualArgs, tt.expectedArgs) {
				t.Errorf("parseLineWithQuotes(%q) = %v, want %v", tt.line, actualArgs, tt.expectedArgs)
			}
		})
	}
}
