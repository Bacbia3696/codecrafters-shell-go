package main

import (
	"bufio"
	"fmt"
	"os"
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
			exitCode, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stdout, "Error reading exit code: %v\n", err)
				continue
			}
			os.Exit(exitCode)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			arg1 := args[0]
			if slices.Contains(cmds, arg1) {
				fmt.Println(arg1 + " is a shell builtin")
			} else {
				fmt.Println(arg1 + ": not found")
			}
		default:
			fmt.Println(command + ": command not found")
		}
	}
}
