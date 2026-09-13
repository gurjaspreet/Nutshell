package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/gurjaspreet/Nutshell/internal/shell"
)

func main() {
	shellApp := shell.New()

	completer := shell.NewWordCompleter()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing shell:", err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if line != "" && shellApp.ExecuteLine(line) {
			return
		}

		if err != nil {
			if err != io.EOF && err != readline.ErrInterrupt {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
			return
		}
	}
}
