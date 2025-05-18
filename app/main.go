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

		commandList, outputFile, err := parser.ParseLine(inputLine)
		// Remove or comment out the debug print for outputFile
		// fmt.Println("outputFile: ", outputFile)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing command: %v\n", err)
			continue
		}

		if len(commandList) == 0 {
			continue
		}

		// Execute command in an anonymous function to scope the defer correctly
		func() {
			originalStdout := os.Stdout
			var file *os.File

			if outputFile != "" {
				var createErr error
				// Create or truncate the file. Open with write-only, create if not exists, truncate if exists.
				file, createErr = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if createErr != nil {
					fmt.Fprintf(os.Stderr, "Error opening output file %s: %v\n", outputFile, createErr)
					return // Return from anonymous function, effectively a continue for the loop
				}
				os.Stdout = file // Redirect stdout to the file

				defer func() {
					os.Stdout = originalStdout // Restore original stdout
					if file != nil {
						file.Close()
					}
				}()
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
		}() // End of anonymous function call
	}
}
