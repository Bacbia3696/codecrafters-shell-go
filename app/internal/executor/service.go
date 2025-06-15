package executor

import (
	"io"
)

// Service provides command execution functionality
type Service struct{}

// NewService creates a new executor service
func NewService() *Service {
	return &Service{}
}

// Execute executes an external command with the provided IO streams
func (s *Service) Execute(command string, args []string, stdin io.Reader, stdout, stderr io.Writer) {
	HandleExternalCommandWithIO(command, args, stdin, stdout, stderr)
}

// FindCommand finds the full path of a command in the system PATH
func (s *Service) FindCommand(commandName string) string {
	return GetCommand(commandName)
}
