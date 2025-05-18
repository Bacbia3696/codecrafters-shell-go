package shell

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

// testShell creates a shell instance with custom IO for testing
func testShell() (*Shell, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	shell := &Shell{
		stdin:         inBuf,
		stdout:        outBuf,
		stderr:        errBuf,
		prompt:        "$ ",
		builtinCmds:   make(map[string]CommandHandler),
		commandFinder: func(string) string { return "" }, // Mock finder that finds nothing
	}
	shell.reader = bufio.NewReader(strings.NewReader(""))

	// Register test command handlers
	shell.builtinCmds["echo"] = func(args []string) {
		outBuf.WriteString(strings.Join(args, " "))
		outBuf.WriteString("\n")
	}

	shell.builtinCmds["fail"] = func(args []string) {
		errBuf.WriteString("Command failed: " + strings.Join(args, " "))
		errBuf.WriteString("\n")
	}

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
