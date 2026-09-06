package shell

import (
	"fmt"
	"os"
	"os/exec"
)

func execute(path string, args []string, redirection Redirection) {
	command := exec.Command(path, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	file, err := openRedirect(redirection)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating output file:", err)
		return
	}
	if file != nil {
		defer file.Close()
		if redirection.FileDescriptor == 1 {
			command.Stdout = file
		} else {
			command.Stderr = file
		}
	}

	if err := command.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
