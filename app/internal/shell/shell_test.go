package shell

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/internal/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/internal/executor"
	"github.com/codecrafters-io/shell-starter-go/app/internal/parser"
)

// testShell creates a shell instance with custom IO for testing
func testShell() (*Shell, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	shell := &Shell{
		stdin:     inBuf,
		stdout:    outBuf,
		stderr:    errBuf,
		prompt:    "$ ",
		builtins:  builtins.NewRegistry(outBuf, errBuf),
		ioManager: NewIOManager(outBuf, errBuf),
		executor:  executor.NewService(),
		parser:    parser.NewService(),
	}
	shell.reader = bufio.NewReader(strings.NewReader(""))

	// Configure with a mock command finder that finds nothing
	shell.builtins.SetCommandFinder(func(string) string { return "" })

	// Register a test command that writes to stderr for testing error handling
	shell.builtins.Register("fail", func(args []string, stdout, stderr io.Writer) error {
		stderr.Write([]byte("Command failed: " + strings.Join(args, " ") + "\n"))
		return nil
	})

	return shell, inBuf, outBuf, errBuf
}

func TestShellExecuteBuiltin(t *testing.T) {
	shell, _, outBuf, errBuf := testShell()

	shell.Execute("echo hello world")

	if outBuf.String() != "hello world\n" {
		t.Errorf("Expected output to be 'hello world\\n', but got %q", outBuf.String())
	}

	if errBuf.String() != "" {
		t.Errorf("Expected no errors, but got: %q", errBuf.String())
	}
}

func TestShellExecuteBuiltinError(t *testing.T) {
	shell, _, outBuf, errBuf := testShell()

	shell.Execute("fail this command")

	if outBuf.String() != "" {
		t.Errorf("Expected no output, but got %q", outBuf.String())
	}

	expected := "Command failed: this command\n"
	if errBuf.String() != expected {
		t.Errorf("Expected error output %q, but got: %q", expected, errBuf.String())
	}
}

func TestShellExecuteExternalNotFound(t *testing.T) {
	shell, _, outBuf, errBuf := testShell()

	// This uses our mock command finder that always returns ""
	shell.Execute("nonexistent arg1 arg2")

	if outBuf.String() != "" {
		t.Errorf("Expected no output, but got %q", outBuf.String())
	}

	if !strings.Contains(errBuf.String(), "command not found") {
		t.Errorf("Expected 'command not found' error, but got: %q", errBuf.String())
	}
}

func TestShellIsBuiltin(t *testing.T) {
	shell, _, _, _ := testShell()

	if !shell.IsBuiltin("echo") {
		t.Error("Expected 'echo' to be recognized as builtin")
	}

	if shell.IsBuiltin("nonexistent") {
		t.Error("Expected 'nonexistent' to not be recognized as builtin")
	}
}

func TestShellExecuteWithRedirection(t *testing.T) {
	shell, _, outBuf, errBuf := testShell()

	// Test stdout redirection - this should not appear in outBuf since it's redirected
	shell.Execute("echo hello > /tmp/test_output.txt")

	// The output should be empty since it was redirected
	if outBuf.String() != "" {
		t.Errorf("Expected no output to stdout, but got %q", outBuf.String())
	}

	if errBuf.String() != "" {
		t.Errorf("Expected no errors, but got: %q", errBuf.String())
	}
}

func TestShellExecuteEmptyCommand(t *testing.T) {
	shell, _, outBuf, errBuf := testShell()

	shell.Execute("")

	if outBuf.String() != "" {
		t.Errorf("Expected no output for empty command, but got %q", outBuf.String())
	}

	if errBuf.String() != "" {
		t.Errorf("Expected no errors for empty command, but got: %q", errBuf.String())
	}
}
