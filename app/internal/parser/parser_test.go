package parser

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name               string
		line               string
		expectedArgs       []string
		expectedOutputFile string
		expectedErrorFile  string
		expectError        bool
	}{
		{
			name:               "empty string",
			line:               "",
			expectedArgs:       []string{},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "only spaces",
			line:               "   ",
			expectedArgs:       []string{},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "simple command",
			line:               "echo",
			expectedArgs:       []string{"echo"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "command with one arg",
			line:               "echo hello",
			expectedArgs:       []string{"echo", "hello"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "command with multiple args",
			line:               "ls -l /tmp",
			expectedArgs:       []string{"ls", "-l", "/tmp"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "command with single quoted arg",
			line:               "echo 'hello world'",
			expectedArgs:       []string{"echo", "hello world"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "command with double quoted arg",
			line:               "echo \"hello world\"",
			expectedArgs:       []string{"echo", "hello world"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "simple output redirection >",
			line:               "ls > out.txt",
			expectedArgs:       []string{"ls"},
			expectedOutputFile: "out.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "output redirection 1>",
			line:               "ls 1> out.txt",
			expectedArgs:       []string{"ls"},
			expectedOutputFile: "out.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "error redirection 2>",
			line:               "ls 2> err.txt",
			expectedArgs:       []string{"ls"},
			expectedOutputFile: "",
			expectedErrorFile:  "err.txt",
			expectError:        false,
		},
		{
			name:               "output redirection with args >",
			line:               "echo hello world > message.txt",
			expectedArgs:       []string{"echo", "hello", "world"},
			expectedOutputFile: "message.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "error redirection with args 2>",
			line:               "echo hello world 2> errors.log",
			expectedArgs:       []string{"echo", "hello", "world"},
			expectedOutputFile: "",
			expectedErrorFile:  "errors.log",
			expectError:        false,
		},
		{
			name:               "redirection with no space before >",
			line:               "ls>out.txt",
			expectedArgs:       []string{"ls"},
			expectedOutputFile: "out.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "redirection with no space before 2>",
			line:               "ls2>err.txt",
			expectedArgs:       []string{"ls2>err.txt"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "redirection 2> with space before but not part of command",
			line:               "ls 2> err.txt",
			expectedArgs:       []string{"ls"},
			expectedOutputFile: "",
			expectedErrorFile:  "err.txt",
			expectError:        false,
		},
		{
			name:               "redirection 2> with space before and command arg",
			line:               "ls -1 file 2> err.txt",
			expectedArgs:       []string{"ls", "-1", "file"},
			expectedOutputFile: "",
			expectedErrorFile:  "err.txt",
			expectError:        false,
		},
		{
			name:               "redirection with quoted filename >",
			line:               "echo test > \"file with spaces.txt\"",
			expectedArgs:       []string{"echo", "test"},
			expectedOutputFile: "file with spaces.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "redirection with quoted filename ' >",
			line:               "echo test > 'file with spaces.txt'",
			expectedArgs:       []string{"echo", "test"},
			expectedOutputFile: "file with spaces.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "redirection with quoted filename 2>",
			line:               "echo test 2> \"errors file.log\"",
			expectedArgs:       []string{"echo", "test"},
			expectedOutputFile: "",
			expectedErrorFile:  "errors file.log",
			expectError:        false,
		},
		{
			name:               "redirection operator inside quotes (should not redirect)",
			line:               "echo \"hello > world\" foo",
			expectedArgs:       []string{"echo", "hello > world", "foo"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "2> operator inside quotes (should not redirect)",
			line:               "echo \"hello 2> world\" foo",
			expectedArgs:       []string{"echo", "hello 2> world", "foo"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "missing filename after >",
			line:               "echo hello >",
			expectedArgs:       nil,
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        true,
		},
		{
			name:               "missing filename after 1>",
			line:               "echo hello 1>",
			expectedArgs:       nil,
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        true,
		},
		{
			name:               "missing filename after 2>",
			line:               "echo hello 2>",
			expectedArgs:       nil,
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        true,
		},
		{
			name:               "only redirection operator >",
			line:               "> out.txt",
			expectedArgs:       []string{},
			expectedOutputFile: "out.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "only redirection operator 2>",
			line:               "2> err.txt",
			expectedArgs:       []string{},
			expectedOutputFile: "",
			expectedErrorFile:  "err.txt",
			expectError:        false,
		},
		{
			name:               "Filename with special characters, unquoted",
			line:               "cmd > file-name_with.chars.123",
			expectedArgs:       []string{"cmd"},
			expectedOutputFile: "file-name_with.chars.123",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "Filename that looks like a redirection but isn't (part of arg)",
			line:               "echo arg1>notfile arg2",
			expectedArgs:       []string{"echo", "arg1>notfile", "arg2"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "command with multiple single quoted args",
			line:               "echo 'hello world' 'another one'",
			expectedArgs:       []string{"echo", "hello world", "another one"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "command with multiple double quoted args",
			line:               "echo \"hello world\" \"another one\"",
			expectedArgs:       []string{"echo", "hello world", "another one"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "mixed quoted and unquoted",
			line:               "command 'quoted arg' unquoted \"another quoted\" last",
			expectedArgs:       []string{"command", "quoted arg", "unquoted", "another quoted", "last"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "leading and trailing spaces",
			line:               "  echo hello  ",
			expectedArgs:       []string{"echo", "hello"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "multiple spaces between args",
			line:               "echo   hello   world",
			expectedArgs:       []string{"echo", "hello", "world"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "quoted arg with leading/trailing space inside quotes",
			line:               "echo '  spaced arg  '",
			expectedArgs:       []string{"echo", "  spaced arg  "},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "adjacent quoted args",
			line:               "echo 'hello''world'",
			expectedArgs:       []string{"echo", "helloworld"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "adjacent quoted args with double quotes",
			line:               "echo \"hello\"\"world\"",
			expectedArgs:       []string{"echo", "helloworld"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "adjacent mixed quotes",
			line:               "echo 'hello'\"world\"",
			expectedArgs:       []string{"echo", "helloworld"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "complex case with multiple quotes and spaces",
			line:               " cmd 'arg1 part1' '' arg2 'arg3 part1 part2' ",
			expectedArgs:       []string{"cmd", "arg1 part1", "", "arg2", "arg3 part1 part2"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "line ending with an open quote (should be treated as unclosed)",
			line:               "echo 'hello",
			expectedArgs:       []string{"echo", "hello"},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "command with empty single quoted arg",
			line:               "echo ''",
			expectedArgs:       []string{"echo", ""},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "command with empty double quoted arg",
			line:               "echo \"\"",
			expectedArgs:       []string{"echo", ""},
			expectedOutputFile: "",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "output redirection 1> with no space after",
			line:               "ls 1>out.txt",
			expectedArgs:       []string{"ls"},
			expectedOutputFile: "out.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "output redirection 1> with no space before (but space after 1)",
			line:               "ls 1 > out.txt",
			expectedArgs:       []string{"ls"},
			expectedOutputFile: "out.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "output redirection with args 1>",
			line:               "echo hello world 1> message.txt",
			expectedArgs:       []string{"echo", "hello", "world"},
			expectedOutputFile: "message.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
		{
			name:               "redirection with no space after >",
			line:               "ls >out.txt",
			expectedArgs:       []string{"ls"},
			expectedOutputFile: "out.txt",
			expectedErrorFile:  "",
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, outputFile, errorFile, err := ParseLine(tt.line)

			if (err != nil) != tt.expectError {
				t.Errorf("ParseLine() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !reflect.DeepEqual(args, tt.expectedArgs) {
				t.Errorf("ParseLine() args = %v, want %v", args, tt.expectedArgs)
			}
			if outputFile != tt.expectedOutputFile {
				t.Errorf("ParseLine() outputFile = %v, want %v", outputFile, tt.expectedOutputFile)
			}
			if errorFile != tt.expectedErrorFile {
				t.Errorf("ParseLine() errorFile = %v, want %v", errorFile, tt.expectedErrorFile)
			}
		})
	}
}

func TestParseLineEmpty(t *testing.T) {
	args, stdout, stderr, err := ParseLine("")
	if err != nil {
		t.Errorf("Expected no error for empty line but got: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("Expected empty args but got: %v", args)
	}
	if stdout != "" {
		t.Errorf("Expected empty stdout but got: %s", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected empty stderr but got: %s", stderr)
	}
}

func TestParseLineSimpleCommand(t *testing.T) {
	args, stdout, stderr, err := ParseLine("echo hello world")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"echo", "hello", "world"}
	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Expected args %v but got: %v", expected, args)
	}
	if stdout != "" || stderr != "" {
		t.Errorf("Expected empty stdout/stderr but got stdout=%s stderr=%s", stdout, stderr)
	}
}

func TestParseLineWithQuotes(t *testing.T) {
	args, stdout, stderr, err := ParseLine("echo 'hello world' \"another quote\"")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"echo", "hello world", "another quote"}
	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Expected args %v but got: %v", expected, args)
	}
	if stdout != "" || stderr != "" {
		t.Errorf("Expected empty stdout/stderr but got stdout=%s stderr=%s", stdout, stderr)
	}
}

func TestParseLineStdoutRedirection(t *testing.T) {
	args, stdout, stderr, err := ParseLine("echo hello > out.txt")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"echo", "hello"}
	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Expected args %v but got: %v", expected, args)
	}
	if stdout != "out.txt" {
		t.Errorf("Expected stdout=out.txt but got: %s", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected empty stderr but got: %s", stderr)
	}
}

func TestParseLineStderrRedirection(t *testing.T) {
	args, stdout, stderr, err := ParseLine("ls /nonexistent 2> err.txt")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"ls", "/nonexistent"}
	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Expected args %v but got: %v", expected, args)
	}
	if stdout != "" {
		t.Errorf("Expected empty stdout but got: %s", stdout)
	}
	if stderr != "err.txt" {
		t.Errorf("Expected stderr=err.txt but got: %s", stderr)
	}
}

func TestParseLineNestedQuotes(t *testing.T) {
	args, stdout, stderr, err := ParseLine("echo \"It's a 'quoted' string\"")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"echo", "It's a 'quoted' string"}
	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Expected args %v but got: %v", expected, args)
	}
	if stdout != "" || stderr != "" {
		t.Errorf("Expected empty stdout/stderr but got stdout=%s stderr=%s", stdout, stderr)
	}
}

func TestParseLineEmptyQuotes(t *testing.T) {
	args, stdout, stderr, err := ParseLine("echo '' \"\"")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"echo", "", ""}
	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Expected args %v but got: %v", expected, args)
	}
	if stdout != "" || stderr != "" {
		t.Errorf("Expected empty stdout/stderr but got stdout=%s stderr=%s", stdout, stderr)
	}
}

func TestParseLineMissingFilename(t *testing.T) {
	_, _, _, err := ParseLine("echo hello >")
	if err == nil {
		t.Errorf("Expected error for missing filename but got none")
	}
}

func TestParseLineQuotedRedirection(t *testing.T) {
	args, stdout, stderr, err := ParseLine("echo hello > 'my file.txt'")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"echo", "hello"}
	if !reflect.DeepEqual(args, expected) {
		t.Errorf("Expected args %v but got: %v", expected, args)
	}
	if stdout != "my file.txt" {
		t.Errorf("Expected stdout='my file.txt' but got: %s", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected empty stderr but got: %s", stderr)
	}
}

func TestParseLineAppendRedirection(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedArgs   []string
		expectedOutput string
		expectedError  string
		outputAppend   bool
		errorAppend    bool
	}{
		{
			name:           "simple append redirection >>",
			input:          "echo hello >> out.txt",
			expectedArgs:   []string{"echo", "hello"},
			expectedOutput: "out.txt",
			outputAppend:   true,
		},
		{
			name:           "append redirection 1>>",
			input:          "echo hello 1>> out.txt",
			expectedArgs:   []string{"echo", "hello"},
			expectedOutput: "out.txt",
			outputAppend:   true,
		},
		{
			name:          "error append redirection 2>>",
			input:         "ls /nonexistent 2>> err.txt",
			expectedArgs:  []string{"ls", "/nonexistent"},
			expectedError: "err.txt",
			errorAppend:   true,
		},
		{
			name:           "append with spaces around >>",
			input:          "echo test  >>  output.log",
			expectedArgs:   []string{"echo", "test"},
			expectedOutput: "output.log",
			outputAppend:   true,
		},
		{
			name:           "append with quoted filename",
			input:          "echo data >> \"my file.txt\"",
			expectedArgs:   []string{"echo", "data"},
			expectedOutput: "my file.txt",
			outputAppend:   true,
		},
		{
			name:           "append with single quoted filename",
			input:          "echo data >> 'my file.txt'",
			expectedArgs:   []string{"echo", "data"},
			expectedOutput: "my file.txt",
			outputAppend:   true,
		},
		{
			name:         "append operator inside quotes should not redirect",
			input:        "echo 'hello >> world'",
			expectedArgs: []string{"echo", "hello >> world"},
		},
		{
			name:         "2>> operator inside quotes should not redirect",
			input:        "echo \"error 2>> log\"",
			expectedArgs: []string{"echo", "error 2>> log"},
		},
		{
			name:           "mixed regular and append redirection - regular first",
			input:          "echo hello > out.txt",
			expectedArgs:   []string{"echo", "hello"},
			expectedOutput: "out.txt",
			outputAppend:   false,
		},
		{
			name:         "no space before >> (glued to argument)",
			input:        "echo hello>>out.txt",
			expectedArgs: []string{"echo", "hello>>out.txt"},
		},
		{
			name:          "2>> with space before digit",
			input:         "ls /foo 2 >> err.txt",
			expectedArgs:  []string{"ls", "/foo"},
			expectedError: "err.txt",
			errorAppend:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, outputFile, errorFile, outputAppend, errorAppend, err := ParseLineWithMode(tt.input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(args, tt.expectedArgs) {
				t.Errorf("args = %v, want %v", args, tt.expectedArgs)
			}
			if outputFile != tt.expectedOutput {
				t.Errorf("outputFile = %q, want %q", outputFile, tt.expectedOutput)
			}
			if errorFile != tt.expectedError {
				t.Errorf("errorFile = %q, want %q", errorFile, tt.expectedError)
			}
			if outputAppend != tt.outputAppend {
				t.Errorf("outputAppend = %v, want %v", outputAppend, tt.outputAppend)
			}
			if errorAppend != tt.errorAppend {
				t.Errorf("errorAppend = %v, want %v", errorAppend, tt.errorAppend)
			}
		})
	}
}
