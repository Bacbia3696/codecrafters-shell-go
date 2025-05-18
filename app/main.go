package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/internal/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/internal/executor"
	"github.com/codecrafters-io/shell-starter-go/app/internal/parser"
)

func main() {
	// Configure the command finder for the builtins package
	// This allows builtins.HandleType to use executor.GetCommand without a direct import cycle.
	builtins.ConfigureCommandFinder(executor.GetCommand)

	for {
		fmt.Fprint(os.Stdout, "$ ")
		inputLine, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" { // Check for EOF to exit gracefully
				fmt.Fprintln(os.Stdout, "exit")
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "Error reading command: %v\n", err)
			continue // Continue to the next prompt on read error
		}

		commandList := parser.ParseLine(inputLine)

		if len(commandList) == 0 {
			continue
		}

		command := commandList[0]
		args := commandList[1:]

		switch command {
		case "exit":
			builtins.HandleExit(args)
		case "echo":
			builtins.HandleEcho(args)
		case "pwd":
			builtins.HandlePwd()
		case "cd":
			builtins.HandleCd(args)
		case "type":
			builtins.HandleType(args)
		default:
			executor.HandleExternalCommand(command, args)
		}
	}
}
