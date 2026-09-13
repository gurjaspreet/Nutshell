package shell

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var supportedCommands = []string{"cd", "echo", "exit", "pwd", "type", "refresh"}

func runBuiltin(command string, args []string, redirection Redirection) bool {
	switch command {
	case "cd":
		changeDirectory(args)
	case "echo":
		echo(args, redirection)
	case "pwd":
		printWorkingDirectory(redirection)
	case "refresh":
		refreshPathExecutables()
	case "type":
		if len(args) > 0 {
			commandType(args[0], redirection)
		}
	case "exit":
		return true
	default:
		return false
	}
	return true
}

func changeDirectory(args []string) {
	if len(args) == 0 {
		return
	}

	path := args[0]
	if path == "~" {
		homeDirectory, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting home directory:", err)
			return
		}
		path = homeDirectory
	}

	if err := os.Chdir(path); err != nil {
		fmt.Println("cd: " + path + ": No such file or directory")
	}
}

func echo(args []string, redirection Redirection) {
	writeOutput(1, strings.Join(args, " ")+"\n", redirection)
}

func printWorkingDirectory(redirection Redirection) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		writeOutput(2, "Error getting current directory: "+err.Error()+"\n", redirection)
		return
	}
	writeOutput(1, workingDirectory+"\n", redirection)
}

func commandType(command string, redirection Redirection) {
	var output string
	if slices.Contains(supportedCommands, command) {
		output = command + " is a shell builtin\n"
	} else if path, err := exec.LookPath(command); err == nil {
		output = command + " is " + path + "\n"
	} else {
		output = command + ": not found\n"
	}
	writeOutput(1, output, redirection)
}
