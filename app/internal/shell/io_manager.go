package shell

import (
	"io"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/internal/errors"
)

// IOManagerImpl manages input/output redirection for the shell
type IOManagerImpl struct {
	originalStdout io.Writer
	originalStderr io.Writer
	currentStdout  io.Writer
	currentStderr  io.Writer
}

// NewIOManager creates a new IO manager
func NewIOManager(stdout, stderr io.Writer) *IOManagerImpl {
	return &IOManagerImpl{
		originalStdout: stdout,
		originalStderr: stderr,
		currentStdout:  stdout,
		currentStderr:  stderr,
	}
}

// SetupRedirection sets up file redirection and returns a cleanup function
func (m *IOManagerImpl) SetupRedirection(outputFile, errorFile string) (cleanup func(), err error) {
	var outFile, errFile *os.File

	// Setup cleanup function that will restore original streams
	cleanup = func() {
		m.currentStdout = m.originalStdout
		m.currentStderr = m.originalStderr
		if outFile != nil {
			outFile.Close()
		}
		if errFile != nil {
			errFile.Close()
		}
	}

	// Setup output redirection
	if outputFile != "" {
		outFile, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			cleanup()
			return nil, errors.NewIOError("opening", outputFile, err.Error())
		}
		m.currentStdout = outFile
	}

	// Setup error redirection
	if errorFile != "" {
		errFile, err = os.OpenFile(errorFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			cleanup()
			return nil, errors.NewIOError("opening", errorFile, err.Error())
		}
		m.currentStderr = errFile
	}

	return cleanup, nil
}

// GetCurrentStreams returns the current stdout and stderr streams
func (m *IOManagerImpl) GetCurrentStreams() (stdout, stderr io.Writer) {
	return m.currentStdout, m.currentStderr
}
