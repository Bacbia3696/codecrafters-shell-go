package builtins

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// CommandHandler defines a function that handles a shell command
type CommandHandler func(args []string, stdout, stderr io.Writer) error

// Registry manages built-in commands
type Registry struct {
	commands      map[string]CommandHandler
	commandFinder func(string) string
}

// NewRegistry creates a new built-in command registry
func NewRegistry(stdout, stderr io.Writer) *Registry {
	r := &Registry{
		commands: make(map[string]CommandHandler),
	}
	r.registerDefaults()
	return r
}

// SetCommandFinder sets the function used to find external commands
func (r *Registry) SetCommandFinder(finder func(string) string) {
	r.commandFinder = finder
}

// IsBuiltin checks if a command is a built-in command
func (r *Registry) IsBuiltin(cmd string) bool {
	_, exists := r.commands[cmd]
	return exists
}

// Execute executes a built-in command with the provided streams
func (r *Registry) Execute(cmd string, args []string, stdout, stderr io.Writer) error {
	handler, exists := r.commands[cmd]
	if !exists {
		return fmt.Errorf("built-in command not found: %s", cmd)
	}
	return handler(args, stdout, stderr)
}

// Register registers a new built-in command
func (r *Registry) Register(cmd string, handler CommandHandler) {
	r.commands[cmd] = handler
}

// registerDefaults registers the default built-in commands
func (r *Registry) registerDefaults() {
	r.commands["exit"] = r.handleExit
	r.commands["echo"] = r.handleEcho
	r.commands["pwd"] = r.handlePwd
	r.commands["cd"] = r.handleCd
	r.commands["type"] = r.handleType
}

// handleExit handles the 'exit' built-in command
func (r *Registry) handleExit(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		os.Exit(0)
		return nil
	}
	exitCode, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "exit: invalid exit code: %s\n", args[0])
		return err
	}
	os.Exit(exitCode)
	return nil
}

// handleEcho handles the 'echo' built-in command
func (r *Registry) handleEcho(args []string, stdout, stderr io.Writer) error {
	fmt.Fprintln(stdout, strings.Join(args, " "))
	return nil
}

// handlePwd handles the 'pwd' built-in command
func (r *Registry) handlePwd(args []string, stdout, stderr io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "pwd: %v\n", err)
		return err
	}
	fmt.Fprintln(stdout, dir)
	return nil
}

// handleCd handles the 'cd' built-in command
func (r *Registry) handleCd(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "cd: missing argument")
		return fmt.Errorf("missing argument")
	}

	targetDir := args[0]
	if targetDir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(stderr, "cd: error getting home directory: %v\n", err)
			return err
		}
		targetDir = homeDir
	}

	err := os.Chdir(targetDir)
	if err != nil {
		fmt.Fprintf(stderr, "cd: %s: No such file or directory\n", args[0])
		return err
	}
	return nil
}

// handleType handles the 'type' built-in command
func (r *Registry) handleType(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "type: missing argument")
		return fmt.Errorf("missing argument")
	}

	cmdName := args[0]
	if r.IsBuiltin(cmdName) {
		fmt.Fprintln(stdout, cmdName+" is a shell builtin")
		return nil
	}

	if r.commandFinder == nil {
		fmt.Fprintf(stderr, "type: command finder not configured\n")
		fmt.Fprintln(stdout, cmdName+": not found")
		return fmt.Errorf("command finder not configured")
	}

	foundPath := r.commandFinder(cmdName)
	if foundPath != "" {
		fmt.Fprintln(stdout, cmdName+" is "+foundPath)
	} else {
		fmt.Fprintln(stdout, cmdName+": not found")
	}
	return nil
}
