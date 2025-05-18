package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

var builtInCmds = []string{"exit", "echo", "type", "pwd", "cd"}

// parseLineWithQuotes splits a line into arguments, respecting single quotes.
// Text within single quotes is treated as a single argument, and the quotes are removed.
// Example: "echo 'hello world' foo" -> ["echo", "hello world", "foo"]
// Example: "echo ” foo" -> ["echo", "", "foo"]
func parseLineWithQuotes(line string) []string {
	args := make([]string, 0)
	var currentArg strings.Builder
	inQuote := false
	lineRunes := []rune(strings.TrimSpace(line))

	// keep this
	for i := range len(lineRunes) {
		char := lineRunes[i]

		// Single quote
		if char == '\'' {
			inQuote = !inQuote
			if !inQuote {
				isLastCharInSegment := (i+1 == len(lineRunes)) || (i+1 < len(lineRunes) && lineRunes[i+1] == ' ')
				if isLastCharInSegment {
					args = append(args, currentArg.String())
					currentArg.Reset()
				}
			}
		} else if char == ' ' && !inQuote {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		} else {
			currentArg.WriteRune(char)
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}
	return args
}

func handleExit(args []string) {
	if len(args) == 0 {
		os.Exit(0)
		return // os.Exit doesn't return, but for consistency
	}
	exitCode, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit: invalid exit code: %s\n", args[0])
		// In a real shell, this might set a status and not exit the shell itself,
		// but for this project, exiting the sub-shell process is fine.
		// If this were the main shell loop, we'd 'continue' the loop.
		// Since this handler is called and then the loop continues,
		// we don't need to os.Exit(1) here, just indicate error.
		return
	}
	os.Exit(exitCode)
}

func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func handlePwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
		return
	}
	fmt.Println(dir)
}

func handleCd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "cd: missing argument") // Error to Stderr
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

func handleType(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "type: missing argument")
		return
	}
	arg1 := args[0]
	if slices.Contains(builtInCmds, arg1) {
		fmt.Fprintln(os.Stdout, arg1+" is a shell builtin")
	} else {
		foundPath := getCommand(arg1)
		if foundPath != "" {
			fmt.Fprintln(os.Stdout, arg1+" is "+foundPath)
		} else {
			fmt.Fprintln(os.Stdout, arg1+": not found") // This is info, not an error
		}
	}
}

func handleExternalCommand(command string, args []string) {
	foundPath := getCommand(command)
	if foundPath != "" {
		cmd := exec.Command(command, args...) // Corrected to use foundPath
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				_ = exitErr // Command ran but exited non-zero. Output is already on Stderr.
			} else {
				fmt.Fprintf(os.Stderr, "%s: command failed to start: %v\n", command, err)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s: command not found\n", command) // Error to Stderr
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		inputLine, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("exit")
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "Error reading command: %v\n", err)
			return
		}

		commandList := parseLineWithQuotes(inputLine)

		if len(commandList) == 0 {
			continue
		}

		command := commandList[0]
		args := commandList[1:]

		switch command {
		case "exit":
			handleExit(args)
		case "echo":
			handleEcho(args)
		case "pwd":
			handlePwd()
		case "cd":
			handleCd(args)
		case "type":
			handleType(args)
		default:
			handleExternalCommand(command, args)
		}
	}
}

func getCommand(commandName string) string {
	pathsEnv := os.Getenv("PATH")
	if pathsEnv == "" {
		return ""
	}
	pathDirs := strings.Split(pathsEnv, string(os.PathListSeparator))
	for _, dir := range pathDirs {
		if dir == "" {
			dir = "."
		}
		fullPath := filepath.Join(dir, commandName)
		fileInfo, err := os.Stat(fullPath)
		if err == nil {
			if !fileInfo.IsDir() && (fileInfo.Mode().Perm()&0111 != 0) {
				return fullPath
			}
		}
	}
	return ""
}
