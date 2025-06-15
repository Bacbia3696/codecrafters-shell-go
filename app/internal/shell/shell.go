package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/internal/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/internal/errors"
	"github.com/codecrafters-io/shell-starter-go/app/internal/executor"
	"github.com/codecrafters-io/shell-starter-go/app/internal/parser"
)

// Shell represents the shell program state and behavior
type Shell struct {
	reader    *bufio.Reader
	stdin     io.Reader
	stdout    io.Writer
	stderr    io.Writer
	prompt    string
	builtins  BuiltinRegistry
	ioManager IOManager
	executor  CommandExecutor
	parser    CommandParser
}

// NewShell creates a new shell instance with default configuration
func NewShell() *Shell {
	stdout := os.Stdout
	stderr := os.Stderr

	return NewShellWithDependencies(
		os.Stdin,
		stdout,
		stderr,
		parser.NewService(),
		executor.NewService(),
		builtins.NewRegistry(stdout, stderr),
		NewIOManager(stdout, stderr),
	)
}

// NewShellWithDependencies creates a new shell instance with provided dependencies
func NewShellWithDependencies(
	stdin io.Reader,
	stdout, stderr io.Writer,
	parser CommandParser,
	executor CommandExecutor,
	builtins BuiltinRegistry,
	ioManager IOManager,
) *Shell {
	s := &Shell{
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
		prompt:    "$ ",
		builtins:  builtins,
		ioManager: ioManager,
		executor:  executor,
		parser:    parser,
	}
	s.reader = bufio.NewReader(s.stdin)

	// Configure the command finder for builtins
	builtins.SetCommandFinder(executor.FindCommand)

	return s
}

// IsBuiltin checks if a command is a built-in command
func (s *Shell) IsBuiltin(cmd string) bool {
	return s.builtins.IsBuiltin(cmd)
}

// Execute executes a single command line
func (s *Shell) Execute(inputLine string) {
	// Try to use the new append-aware parser if available
	if parserWithMode, ok := s.parser.(CommandParserWithMode); ok {
		args, outputFile, errorFile, outputAppend, errorAppend, err := parserWithMode.ParseLineWithMode(inputLine)

		if err != nil {
			fmt.Fprintf(s.stderr, "%s\n", err.Error())
			return
		}

		if len(args) == 0 {
			return // Empty command, nothing to do
		}

		// Setup redirection using IOManager with append mode support
		if ioManagerWithMode, ok := s.ioManager.(IOManagerWithMode); ok {
			cleanup, err := ioManagerWithMode.SetupRedirectionWithMode(outputFile, errorFile, outputAppend, errorAppend)
			if err != nil {
				fmt.Fprintf(s.stderr, "%s\n", err.Error())
				return
			}
			defer cleanup()
		} else {
			// Fallback to basic redirection (no append support)
			cleanup, err := s.ioManager.SetupRedirection(outputFile, errorFile)
			if err != nil {
				fmt.Fprintf(s.stderr, "%s\n", err.Error())
				return
			}
			defer cleanup()
		}

		// Get current streams from IOManager
		currentStdout, currentStderr := s.ioManager.GetCurrentStreams()

		// Get command and arguments
		command := args[0]
		cmdArgs := args[1:]

		// Execute command
		if s.builtins.IsBuiltin(command) {
			err := s.builtins.Execute(command, cmdArgs, currentStdout, currentStderr)
			if err != nil {
				fmt.Fprintf(currentStderr, "%s\n", err.Error())
			}
		} else {
			s.executor.Execute(command, cmdArgs, s.stdin, currentStdout, currentStderr)
		}
	} else {
		// Fallback to original parsing (no append support)
		args, outputFile, errorFile, err := s.parser.ParseLine(inputLine)

		if err != nil {
			fmt.Fprintf(s.stderr, "%s\n", err.Error())
			return
		}

		if len(args) == 0 {
			return // Empty command, nothing to do
		}

		// Setup redirection using IOManager
		cleanup, err := s.ioManager.SetupRedirection(outputFile, errorFile)
		if err != nil {
			fmt.Fprintf(s.stderr, "%s\n", err.Error())
			return
		}
		defer cleanup()

		// Get current streams from IOManager
		currentStdout, currentStderr := s.ioManager.GetCurrentStreams()

		// Get command and arguments
		command := args[0]
		cmdArgs := args[1:]

		// Execute command
		if s.builtins.IsBuiltin(command) {
			err := s.builtins.Execute(command, cmdArgs, currentStdout, currentStderr)
			if err != nil {
				fmt.Fprintf(currentStderr, "%s\n", err.Error())
			}
		} else {
			s.executor.Execute(command, cmdArgs, s.stdin, currentStdout, currentStderr)
		}
	}
}

// Run starts the shell's read-eval-print loop
func (s *Shell) Run() {
	for {
		fmt.Fprint(s.stdout, s.prompt)
		inputLine, err := s.reader.ReadString('\n')

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Fprintln(s.stdout, "exit")
				os.Exit(0)
			}
			ioErr := errors.NewIOError("reading", "command", err.Error())
			fmt.Fprintf(s.stderr, "%s\n", ioErr.Error())
			continue
		}

		// Trim trailing newline for consistent parsing
		inputLine = strings.TrimSuffix(inputLine, "\n")

		s.Execute(inputLine)
	}
}
