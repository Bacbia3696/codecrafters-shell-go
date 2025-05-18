[![progress-banner](https://backend.codecrafters.io/progress/shell/0add5fa9-5180-4f99-a34e-3008551413ec)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This repository contains a Go implementation for the
["Build Your Own Shell" Challenge](https://app.codecrafters.io/courses/shell/overview)
on CodeCrafters.

The goal is to build a POSIX-compliant shell capable of interpreting commands,
running external programs, and handling various shell features.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

## Features

Currently, the shell supports:
*   A Read-Eval-Print Loop (REPL) for interactive command input.
*   Execution of external commands found in the system's PATH.
*   Built-in commands:
    *   `exit [code]` - Exits the shell.
    *   `echo [args...]` - Prints arguments to standard output.
    *   `pwd` - Prints the current working directory.
    *   `cd <directory>` - Changes the current working directory.
    *   `type <command>` - Displays information about a command (builtin or external).
*   Output redirection to a file using `> filename`.
*   Graceful exit on `EOF` (Ctrl+D).

## Prerequisites

*   Go (version 1.24 or later) installed locally.

## Building and Running

1.  Ensure you have Go 1.24 installed.
2.  Run the script `./your_program.sh`. This script will:
    *   Compile the Go source files (from the `app/` directory).
    *   Run the compiled shell program.

Once the shell starts, you'll see a `$` prompt. Type your commands and press Enter.
To exit the shell, you can type `exit` or press `Ctrl+D`.

## Development (CodeCrafters Stages)

This section is relevant for submitting your solution to CodeCrafters.

### Passing the first stage

The entry point for your `shell` implementation is in `app/main.go`. Study and
uncomment the relevant code, and push your changes to pass the first stage:

```sh
git commit -am "pass 1st stage" # any msg
git push origin master
```

### Stage 2 & beyond

1.  Make your changes to implement new features or fix bugs.
2.  Commit your changes: `git commit -am "implemented feature X"`
3.  Push to CodeCrafters: `git push origin master`

Test output will be streamed to your terminal.
