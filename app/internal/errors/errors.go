package errors

import (
	"fmt"
)

// ShellError is the base interface for all shell-specific errors
type ShellError interface {
	error
	ShellError() string
}

// CommandNotFoundError represents an error when a command cannot be found
type CommandNotFoundError struct {
	Command string
}

func (e CommandNotFoundError) Error() string {
	return fmt.Sprintf("%s: command not found", e.Command)
}

func (e CommandNotFoundError) ShellError() string {
	return "command_not_found"
}

// CommandFailedError represents an error when a command fails to execute
type CommandFailedError struct {
	Command string
	Reason  string
}

func (e CommandFailedError) Error() string {
	return fmt.Sprintf("%s: %s", e.Command, e.Reason)
}

func (e CommandFailedError) ShellError() string {
	return "command_failed"
}

// ParseError represents an error during command parsing
type ParseError struct {
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parse error: %s", e.Message)
}

func (e ParseError) ShellError() string {
	return "parse_error"
}

// IOError represents an error during I/O operations
type IOError struct {
	Operation string
	Target    string
	Reason    string
}

func (e IOError) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Operation, e.Target, e.Reason)
}

func (e IOError) ShellError() string {
	return "io_error"
}

// NewParseError creates a new parse error
func NewParseError(msg string) ParseError {
	return ParseError{Message: msg}
}

// NewIOError creates a new IO error
func NewIOError(op, target, reason string) IOError {
	return IOError{
		Operation: op,
		Target:    target,
		Reason:    reason,
	}
}

// NewCommandNotFoundError creates a new command not found error
func NewCommandNotFoundError(cmd string) CommandNotFoundError {
	return CommandNotFoundError{Command: cmd}
}

// NewCommandFailedError creates a new command failed error
func NewCommandFailedError(cmd, reason string) CommandFailedError {
	return CommandFailedError{
		Command: cmd,
		Reason:  reason,
	}
}
