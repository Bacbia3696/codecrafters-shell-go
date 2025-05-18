package builtins

import (
	"fmt"
	"os"
	"slices" // Go 1.21+ for slices package, ensure compatibility or use local equivalent for older versions
	"strconv"
	"strings"
)

// commandFinder is a function signature for a helper that can find commands.
// This will be provided by an external package (e.g., executor).
var commandFinder func(commandName string) string

// ConfigureCommandFinder sets the function used by HandleType to find external commands.
// This should be called once during initialization (e.g., in main).
func ConfigureCommandFinder(finder func(commandName string) string) {
	commandFinder = finder
}

var builtInCmdsList = []string{"exit", "echo", "type", "pwd", "cd"}

// IsBuiltin checks if a command name is a known built-in command.
func IsBuiltin(cmdName string) bool {
	// For older Go versions that don't have slices.Contains:
	// for _, b := range builtInCmdsList {
	// 	if b == cmdName {
	// 		return true
	// 	}
	// }
	// return false
	return slices.Contains(builtInCmdsList, cmdName)
}

// HandleExit handles the 'exit' built-in command.
func HandleExit(args []string) {
	if len(args) == 0 {
		os.Exit(0)
		return
	}
	exitCode, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit: invalid exit code: %s\n", args[0])
		return
	}
	os.Exit(exitCode)
}

// HandleEcho handles the 'echo' built-in command.
func HandleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// HandlePwd handles the 'pwd' built-in command.
func HandlePwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
		return
	}
	fmt.Println(dir)
}

// HandleCd handles the 'cd' built-in command.
func HandleCd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "cd: missing argument")
		return
	}

	targetDir := args[0]
	if targetDir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: error getting home directory: %v\n", err)
			return
		}
		targetDir = homeDir
	}
	err := os.Chdir(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", args[0])
	}
}

// HandleType handles the 'type' built-in command.
func HandleType(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "type: missing argument")
		return
	}
	cmdName := args[0]
	if IsBuiltin(cmdName) {
		fmt.Fprintln(os.Stdout, cmdName+" is a shell builtin")
	} else {
		if commandFinder == nil {
			fmt.Fprintf(os.Stderr, "type: command finder not configured\n")
			fmt.Fprintln(os.Stdout, cmdName+": not found")
			return
		}
		foundPath := commandFinder(cmdName)
		if foundPath != "" {
			fmt.Fprintln(os.Stdout, cmdName+" is "+foundPath)
		} else {
			fmt.Fprintln(os.Stdout, cmdName+": not found")
		}
	}
}
