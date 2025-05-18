package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/internal/errors"
)

// GetCommand finds the full path of an executable command in the PATH.
func GetCommand(commandName string) string {
	pathsEnv := os.Getenv("PATH")
	if pathsEnv == "" {
		return ""
	}
	// Use strings.Split instead of SplitSeq for standard library compatibility if SplitSeq is not available.
	// os.PathListSeparator is the correct way to split PATH.
	for _, dir := range strings.Split(pathsEnv, string(os.PathListSeparator)) {
		if dir == "" {
			// This case can happen if PATH starts or ends with a separator, or has adjacent separators.
			// In many shells, an empty directory in PATH (e.g. ":/bin" or "/bin::/usr/bin") implies current directory.
			// However, for security, explicitly using "." is better if that's the desired behavior.
			// For now, let's skip empty dirs to avoid ambiguity, unless current dir search is explicit.
			// The original code used dir = ".", let's stick to that for now if it was intentional for current dir search.
			dir = "." // Match original behavior
		}
		fullPath := filepath.Join(dir, commandName)
		fileInfo, err := os.Stat(fullPath)
		if err == nil {
			if !fileInfo.IsDir() && (fileInfo.Mode().Perm()&0111 != 0) { // Check if executable
				return fullPath
			}
		}
	}
	return ""
}

// HandleExternalCommand executes an external command.
func HandleExternalCommand(command string, args []string) {
	HandleExternalCommandWithIO(command, args, os.Stdin, os.Stdout, os.Stderr)
}

// HandleExternalCommandWithIO executes an external command with custom IO streams.
func HandleExternalCommandWithIO(command string, args []string, stdin io.Reader, stdout, stderr io.Writer) {
	foundPath := GetCommand(command)
	if foundPath != "" {
		// keep using command, because test case assert command is first argument
		cmd := exec.Command(command, args...)
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				// Command ran but exited non-zero. Output is already on Stderr by the command itself.
				// We don't need to print an additional error message for this case.
				_ = exitErr // Suppress unused variable warning
			} else {
				// This error means the command failed to start
				shellErr := errors.NewCommandFailedError(command, err.Error())
				fmt.Fprintf(stderr, "%s\n", shellErr.Error())
			}
		}
	} else {
		shellErr := errors.NewCommandNotFoundError(command)
		fmt.Fprintf(stderr, "%s\n", shellErr.Error())
	}
}
