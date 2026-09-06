package shell

import "strings"

type Lexer struct{}

func (Lexer) Tokenize(command string) []string {
	commandParts := make([]string, 0)
	var part strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	flush := func() {
		if part.Len() > 0 {
			commandParts = append(commandParts, part.String())
			part.Reset()
		}
	}

	for i := 0; i < len(command); i++ {
		ch := command[i]

		if !inDoubleQuote && ch == '\'' {
			inSingleQuote = !inSingleQuote
			continue
		}

		if !inSingleQuote && ch == '"' {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if !inSingleQuote && !inDoubleQuote && (ch == ' ' || ch == '\t') {
			flush()
			continue
		}

		if !inSingleQuote && ch == '\\' {
			i++
			if i < len(command) {
				part.WriteByte(command[i])
			}
			continue
		}

		part.WriteByte(ch)
	}

	flush()
	return commandParts
}
