package shell

import (
	"fmt"
	"os/exec"
	"strings"
)

type Shell struct {
	lexer Lexer
}

func New() *Shell {
	return &Shell{}
}

func (s *Shell) ExecuteLine(command string) bool {
	command = strings.TrimSpace(command)
	if command == "" {
		return false
	}

	commandParts := s.lexer.Tokenize(command)
	commandArgs, redirection := ParseRedirection(commandParts)
	if len(commandArgs) == 0 {
		return false
	}

	if runBuiltin(commandArgs[0], commandArgs[1:], redirection) {
		return commandArgs[0] == "exit"
	}

	if _, err := exec.LookPath(commandArgs[0]); err != nil {
		fmt.Println(command + ": command not found")
		return false
	}

	execute(commandArgs[0], commandArgs[1:], redirection)
	return false
}
