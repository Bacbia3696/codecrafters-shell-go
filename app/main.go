package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

var cmds = []string{"exit", "echo", "type"}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		commands, err := bufio.NewReader(os.Stdin).ReadString('\n')
		commandList := strings.Split(strings.TrimSpace(commands), " ")
		command := commandList[0]
		args := commandList[1:]
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error reading command: %v\n", err)
			return
		}
		switch command {
		case "exit":
			if len(args) == 0 {
				os.Exit(0)
				continue
			}
			exitCode, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stdout, "Error reading exit code: %v\n", err)
				continue
			}
			os.Exit(exitCode)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			if len(args) == 0 {
				fmt.Fprintln(os.Stdout, "type: missing argument")
				continue
			}
			arg1 := args[0]
			if slices.Contains(cmds, arg1) {
				fmt.Println(arg1 + " is a shell builtin")
			} else {
				foundPath := getCommand(arg1)
				if foundPath != "" {
					fmt.Println(arg1 + " is " + foundPath)
				} else {
					fmt.Println(arg1 + ": not found")
				}
			}
		default:
			fmt.Println(command + ": command not found")
		}
	}
}

func getCommand(commandName string) string {
	pathsEnv := os.Getenv("PATH")
	if pathsEnv == "" {
		return ""
	}
	pathDirs := strings.SplitSeq(pathsEnv, string(os.PathListSeparator))
	for dir := range pathDirs {
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
