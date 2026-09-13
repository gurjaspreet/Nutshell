package shell

import (
	"fmt"
	"strings"
)

type wordCompleter struct {
	commands            []string
	includePathCommands bool
}

func NewWordCompleter() *wordCompleter {
	return &wordCompleter{
		commands:            supportedCommands,
		includePathCommands: true,
	}
}

func (c *wordCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	// Only get what's before the cursor
	lineToCursor := string(line[:pos])

	// Get the last word and whether it's a command or an argument
	lastSpace := strings.LastIndexByte(lineToCursor, ' ')
	lastWord := lineToCursor[lastSpace+1:]

	var candidates []string
	if lastSpace == -1 {
		// If there's no space, we're completing a command
		candidates = c.commands
		if c.includePathCommands {
			candidates = append(append([]string(nil), candidates...), pathExecutables...)
		}
	} else {
		// If there's a space, we're completing an argument
		candidates = c.completeArgument(lastWord)
	}

	// Collect full names of everything matching the last word
	var matches []string
	for _, cand := range candidates {
		if strings.HasPrefix(cand, lastWord) {
			matches = append(matches, cand)
		}
	}

	length = len(lastWord)

	// If there's a shared prefix longer than what's typed, complete to that
	// instead of listing every match.
	if len(matches) > 1 {
		lcp := longestCommonPrefix(matches)
		if len(lcp) > len(lastWord) {
			newLine = append(newLine, []rune(lcp[len(lastWord):]))
			return newLine, length
		}
	}

	// Prefix-match the candidates with the last word
	for _, cand := range matches {
		suffix := cand[len(lastWord):]
		newLine = append(newLine, []rune(suffix+" "))
	}

	if len(newLine) != 1 {
		fmt.Print("\a") // Beep (Bell character)
	}

	return newLine, length
}

func longestCommonPrefix(strs []string) string {
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}

func (c *wordCompleter) completeArgument(lastWord string) []string {
	return nil
}
