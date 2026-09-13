# Nutshell

A small Unix-style shell written in Go as part of the
[CodeCrafters Build Your Own Shell](https://app.codecrafters.io/courses/shell/overview)
challenge. The project focuses on command-line parsing, shell builtins,
process execution, and file-descriptor redirection.

## Features

- Interactive REPL with a `$` prompt
- Builtin commands: `cd`, `echo`, `exit`, `pwd`, `refresh`, and `type`
- Execution of external programs available on `PATH`
- Command tab completion for builtins and executables on `PATH`
- Single and double-quoted arguments
- Backslash escaping outside single quotes
- Standard output and standard error redirection with `>`, `>>`, `2>`, and `2>>`
- Home-directory expansion for `cd ~`

## Requirements

- Go 1.26 or later

## Run Locally

Start the shell with:

```sh
go run ./cmd/shell
```

Build a standalone executable with:

```sh
go build -o shell ./cmd/shell
```

Example session:

```text
$ pwd
/path/to/Build you own Shell
$ echo "hello shell"
hello shell
$ echo "saved output" > output.txt
$ type pwd
pwd is a shell builtin
$ refresh
$ exit
```

## Project Structure

```text
.
├── cmd/
│   └── shell/
│       └── main.go          # REPL entrypoint
├── internal/
│   └── shell/
│       ├── lexer.go         # Tokenization and quote handling
│       ├── parser.go        # Command redirection parsing
│       ├── shell.go         # Command orchestration
│       ├── builtins.go      # Builtin command implementations
│       ├── executor.go      # External process execution
│       ├── redirect.go      # Output file handling
│       ├── tab_completer.go # Command tab completion
│       ├── shell_test.go    # Lexer and parser tests
│       ├── command_test.go  # Shell command tests
│       └── tab_completer_test.go # Completion tests
├── go.mod
└── README.md
```

## Testing

Run the complete test suite with:

```sh
go test ./...
```
