package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	foundPath := GetCommand(command) // Use the GetCommand from this package
	if foundPath != "" {
		// The original code used `exec.Command(command, args...)` which might rely on PATH search again.
		// It's better to use `foundPath` as the command to execute directly.
		cmd := exec.Command(foundPath, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				// Command ran but exited non-zero. Output is already on Stderr by the command itself.
				// We don't need to print an additional error message for this case.
				_ = exitErr // Suppress unused variable warning
			} else {
				// This error means the command failed to start (e.g., permission issues not caught by GetCommand,
				// or other exec.Command issues).
				fmt.Fprintf(os.Stderr, "%s: command failed to start: %v\n", command, err)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s: command not found\n", command)
	}
}
