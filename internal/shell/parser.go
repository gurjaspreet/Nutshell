package shell

type Redirection struct {
	FileDescriptor int
	Filename       string
	Append         bool
}

func ParseRedirection(commandParts []string) ([]string, Redirection) {
	filtered := make([]string, 0, len(commandParts))
	redirection := Redirection{FileDescriptor: 1}

	for i := 0; i < len(commandParts); i++ {
		part := commandParts[i]
		switch part {
		case ">", "1>", ">>", "1>>":
			if i+1 < len(commandParts) {
				redirection = Redirection{
					FileDescriptor: 1,
					Filename:       commandParts[i+1],
					Append:         part == ">>" || part == "1>>",
				}
				i++
				continue
			}
		case "2>", "2>>":
			if i+1 < len(commandParts) {
				redirection = Redirection{
					FileDescriptor: 2,
					Filename:       commandParts[i+1],
					Append:         part == "2>>",
				}
				i++
				continue
			}
		}
		filtered = append(filtered, part)
	}

	return filtered, redirection
}
