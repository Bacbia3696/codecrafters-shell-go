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

// CommandHandler defines a function that handles a shell command
type CommandHandler func(args []string)

// Shell represents the shell program state and behavior
type Shell struct {
	reader        *bufio.Reader
	stdin         io.Reader
	stdout        io.Writer
	stderr        io.Writer
	prompt        string
	builtinCmds   map[string]CommandHandler
	commandFinder func(string) string
}

// NewShell creates a new shell instance with default configuration
func NewShell() *Shell {
	s := &Shell{
		stdin:         os.Stdin,
		stdout:        os.Stdout,
		stderr:        os.Stderr,
		prompt:        "$ ",
		builtinCmds:   make(map[string]CommandHandler),
		commandFinder: executor.GetCommand,
	}
	s.reader = bufio.NewReader(s.stdin)

	// Register built-in commands
	s.registerBuiltins()

	// Configure the command finder for builtins package
	builtins.ConfigureCommandFinder(s.commandFinder)

	return s
}

// registerBuiltins registers all built-in commands
func (s *Shell) registerBuiltins() {
	s.builtinCmds["exit"] = builtins.HandleExit
	s.builtinCmds["echo"] = builtins.HandleEcho
	s.builtinCmds["pwd"] = func(args []string) { builtins.HandlePwd() }
	s.builtinCmds["cd"] = builtins.HandleCd
	s.builtinCmds["type"] = builtins.HandleType
}

// IsBuiltin checks if a command is a built-in command
func (s *Shell) IsBuiltin(cmd string) bool {
	_, exists := s.builtinCmds[cmd]
	return exists
}

// Execute executes a single command line
func (s *Shell) Execute(inputLine string) {
	commandList, outputFile, errorFile, err := parser.ParseLine(inputLine)

	if err != nil {
		fmt.Fprintf(s.stderr, "%s\n", err.Error())
		return
	}

	if len(commandList) == 0 {
		return // Empty command, nothing to do
	}

	// Handle output redirection
	originalStdout := s.stdout
	originalStderr := s.stderr
	var outFile, errFile *os.File

	// Function to restore original stdout/stderr
	defer func() {
		s.stdout = originalStdout
		s.stderr = originalStderr

		if outFile != nil {
			outFile.Close()
		}
		if errFile != nil {
			errFile.Close()
		}
	}()

	// Set up output redirection if needed
	if outputFile != "" {
		var createErr error
		outFile, createErr = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if createErr != nil {
			ioErr := errors.NewIOError("opening", outputFile, createErr.Error())
			fmt.Fprintf(s.stderr, "%s\n", ioErr.Error())
			return
		}
		s.stdout = outFile
	}

	if errorFile != "" {
		var createErr error
		errFile, createErr = os.OpenFile(errorFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if createErr != nil {
			ioErr := errors.NewIOError("opening", errorFile, createErr.Error())
			fmt.Fprintf(s.stderr, "%s\n", ioErr.Error())
			return
		}
		s.stderr = errFile
	}

	// Get command and arguments
	command := commandList[0]
	args := commandList[1:]

	// Execute command
	if handler, exists := s.builtinCmds[command]; exists {
		handler(args)
	} else {
		s.executeExternal(command, args)
	}
}

// executeExternal executes an external command
func (s *Shell) executeExternal(command string, args []string) {
	// Create a wrapper for executor that uses the shell's stdout/stderr
	executor.HandleExternalCommandWithIO(command, args, s.stdin, s.stdout, s.stderr)
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
