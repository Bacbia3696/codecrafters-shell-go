package shell

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/internal/builtins"
)

// CommandParser defines the interface for parsing command lines
type CommandParser interface {
	ParseLine(line string) (args []string, outputFile string, errorFile string, err error)
}

// CommandExecutor defines the interface for executing external commands
type CommandExecutor interface {
	Execute(command string, args []string, stdin io.Reader, stdout, stderr io.Writer)
	FindCommand(commandName string) string
}

// BuiltinRegistry defines the interface for managing built-in commands
type BuiltinRegistry interface {
	IsBuiltin(cmd string) bool
	Execute(cmd string, args []string, stdout, stderr io.Writer) error
	SetCommandFinder(finder func(string) string)
	Register(cmd string, handler builtins.CommandHandler)
}

// IOManager defines the interface for handling input/output operations
type IOManager interface {
	SetupRedirection(outputFile, errorFile string) (cleanup func(), err error)
	GetCurrentStreams() (stdout, stderr io.Writer)
}
