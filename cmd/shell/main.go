package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/shell"
)

func main() {
	shellApp := shell.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		line, err := reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		if line != "" && shellApp.ExecuteLine(line) {
			return
		}

		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
			return
		}
	}
}
