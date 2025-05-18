package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
		default:
			fmt.Println(command + ": command not found")
		}
	}
}
