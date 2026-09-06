package shell

import (
	"fmt"
	"os"
)

func openRedirect(redirection Redirection) (*os.File, error) {
	if redirection.Filename == "" {
		return nil, nil
	}

	if redirection.Append {
		return os.OpenFile(redirection.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	return os.Create(redirection.Filename)
}

func writeOutput(fileDescriptor int, content string, redirection Redirection) {
	if redirection.Filename == "" || redirection.FileDescriptor != fileDescriptor {
		if fileDescriptor == 2 {
			fmt.Fprint(os.Stderr, content)
		} else {
			fmt.Fprint(os.Stdout, content)
		}
		return
	}

	file, err := openRedirect(redirection)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing to file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing to file:", err)
	}
}
