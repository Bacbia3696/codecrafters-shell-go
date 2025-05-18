package main

import (
	"github.com/codecrafters-io/shell-starter-go/app/internal/shell"
)

func main() {
	// Create and run a new shell instance
	sh := shell.NewShell()
	sh.Run()
}
