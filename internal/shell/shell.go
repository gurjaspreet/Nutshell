package shell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var pathExecutables []string

func init() {
	refreshPathExecutables()
}

func refreshPathExecutables() {
	// Ensure that the PATH environment variable is set
	if os.Getenv("PATH") == "" {
		fmt.Fprintln(os.Stderr, "Warning: PATH environment variable is not set. Some commands may not be found.")
	}

	pathExecutables = getPathExecutables()

	// Debugging: Print the found executables
	// for _, cmd := range pathExecutables {
	// 	fmt.Printf("Found executable: %s\n", cmd)
	// }
}

func getPathExecutables() []string {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil
	}

	paths := strings.Split(pathEnv, string(os.PathListSeparator))
	seen := make(map[string]bool)
	var names []string

	for _, dir := range paths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if seen[name] {
				continue
			}
			if isExecutable(entry) {
				names = append(names, name)
				seen[name] = true
			}
		}
	}

	return names
}

func isExecutable(entry os.DirEntry) bool {
	info, err := entry.Info()
	if err != nil {
		return false
	}

	mode := info.Mode()
	if runtime.GOOS == "windows" {
		return strings.HasSuffix(strings.ToLower(entry.Name()), ".exe") || strings.HasSuffix(strings.ToLower(entry.Name()), ".bat") || strings.HasSuffix(strings.ToLower(entry.Name()), ".cmd")
	}
	return mode&0111 != 0 // Check if any execute bit is set
}

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
