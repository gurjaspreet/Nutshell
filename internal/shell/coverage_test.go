package shell

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestLexerHandlesTabsEscapesAndTrailingBackslash(t *testing.T) {
	got := (Lexer{}).Tokenize("echo\talpha\\ beta 'gamma delta' \\")
	want := []string{"echo", "alpha beta", "gamma delta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokenize() = %#v, want %#v", got, want)
	}
}

func TestParseRedirectionHandlesAllOperatorsAndMissingFilename(t *testing.T) {
	for _, operator := range []string{"1>", ">"} {
		args, redirection := ParseRedirection([]string{"echo", operator, "output.txt"})
		if !reflect.DeepEqual(args, []string{"echo"}) {
			t.Fatalf("ParseRedirection(%q) args = %#v", operator, args)
		}
		if redirection != (Redirection{FileDescriptor: 1, Filename: "output.txt"}) {
			t.Fatalf("ParseRedirection(%q) redirection = %#v", operator, redirection)
		}
	}

	args, redirection := ParseRedirection([]string{"echo", "2>", "errors.log", "1>>", "output.log"})
	if !reflect.DeepEqual(args, []string{"echo"}) {
		t.Fatalf("multiple redirections args = %#v", args)
	}
	if redirection != (Redirection{FileDescriptor: 1, Filename: "output.log", Append: true}) {
		t.Fatalf("multiple redirections = %#v", redirection)
	}

	args, redirection = ParseRedirection([]string{"echo", ">"})
	if !reflect.DeepEqual(args, []string{"echo", ">"}) || redirection != (Redirection{FileDescriptor: 1}) {
		t.Fatalf("missing filename = args %#v, redirection %#v", args, redirection)
	}
}

func TestChangeDirectoryAndPrintWorkingDirectory(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	changeDirectory(nil)
	temporary := t.TempDir()
	t.Cleanup(func() { _ = os.Chdir(original) })
	changeDirectory([]string{temporary})
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if workingDirectory != temporary {
		t.Fatalf("working directory = %q, want %q", workingDirectory, temporary)
	}

	got := captureOutput(t, func() { printWorkingDirectory(Redirection{}) })
	if strings.TrimSpace(got) != temporary {
		t.Fatalf("pwd output = %q, want %q", got, temporary)
	}

	got = captureOutput(t, func() { changeDirectory([]string{filepath.Join(temporary, "missing")}) })
	if !strings.Contains(got, "cd: "+filepath.Join(temporary, "missing")+": No such file or directory") {
		t.Fatalf("invalid cd output = %q", got)
	}
}

func TestChangeDirectoryHome(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	changeDirectory([]string{"~"})
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if workingDirectory != home {
		t.Fatalf("home directory = %q, want %q", workingDirectory, home)
	}
}

func TestExecuteLineHandlesEmptyInputAndBuiltinDispatch(t *testing.T) {
	shellApp := New()
	if shellApp.ExecuteLine("") || shellApp.ExecuteLine("   ") || shellApp.ExecuteLine("> output.txt") {
		t.Fatal("empty command unexpectedly requested exit")
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	temporary := t.TempDir()
	t.Cleanup(func() { _ = os.Chdir(original) })
	if shellApp.ExecuteLine("cd " + temporary) {
		t.Fatal("cd requested exit")
	}
	if shellApp.ExecuteLine("pwd") {
		t.Fatal("pwd requested exit")
	}
	if shellApp.ExecuteLine("refresh") {
		t.Fatal("refresh requested exit")
	}
}

func TestRefreshWarnsWhenPathIsEmpty(t *testing.T) {
	oldPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", ""); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Setenv("PATH", oldPath)
		refreshPathExecutables()
	})

	output := captureOutput(t, refreshPathExecutables)
	if !strings.Contains(output, "Warning: PATH environment variable is not set") {
		t.Fatalf("empty PATH warning = %q", output)
	}
	if pathExecutables != nil {
		t.Fatalf("empty PATH executables = %#v, want nil", pathExecutables)
	}
}

func TestLongestCommonPrefixWithNoSharedCharacters(t *testing.T) {
	if got := longestCommonPrefix([]string{"cat", "echo"}); got != "" {
		t.Fatalf("longestCommonPrefix() = %q, want empty string", got)
	}
}

func TestRedirectionWritesAndAppends(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "output.txt")
	redirection := Redirection{FileDescriptor: 1, Filename: filename}
	writeOutput(1, "first\n", redirection)
	redirection.Append = true
	writeOutput(1, "second\n", redirection)

	contents, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "first\nsecond\n" {
		t.Fatalf("redirected contents = %q", contents)
	}

	got := captureOutput(t, func() {
		writeOutput(2, "stderr\n", Redirection{})
		writeOutput(1, "stdout\n", Redirection{FileDescriptor: 2, Filename: filename})
	})
	if got != "stderr\nstdout\n" {
		t.Fatalf("unmatched descriptor output = %q", got)
	}
}

func TestRedirectionReportsFileErrors(t *testing.T) {
	missingDirectory := filepath.Join(t.TempDir(), "missing", "output.txt")
	got := captureOutput(t, func() {
		writeOutput(1, "content", Redirection{FileDescriptor: 1, Filename: missingDirectory})
		execute("definitely_missing_command_1234", nil, Redirection{FileDescriptor: 1, Filename: missingDirectory})
	})
	if !strings.Contains(got, "Error writing to file:") || !strings.Contains(got, "Error creating output file:") {
		t.Fatalf("redirection errors = %q", got)
	}
}

func TestRefreshUpdatesPathExecutables(t *testing.T) {
	commandName := "refresh_probe"
	commandContents := "#!/bin/sh\nexit 0\n"
	if runtime.GOOS == "windows" {
		commandContents = "@echo off\n"
	}
	commandPath := createTestCommand(t, commandName, commandContents, 0o755)
	oldPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", filepath.Dir(commandPath)); err != nil {
		t.Fatal(err)
	}
	completer := NewWordCompleter()
	t.Cleanup(func() {
		_ = os.Setenv("PATH", oldPath)
		refreshPathExecutables()
	})

	refreshPathExecutables()
	completion, length := completer.Do([]rune(commandName), len(commandName))
	wantSuffix := filepath.Base(commandPath)[len(commandName):] + " "
	if length != len(commandName) || len(completion) != 1 || string(completion[0]) != wantSuffix {
		t.Fatalf("refreshed completion = %#v, length %d; want %q", completion, length, wantSuffix)
	}
}

func TestGetPathExecutablesSkipsInvalidEntriesAndDeduplicates(t *testing.T) {
	first := t.TempDir()
	second := t.TempDir()
	name := "duplicate_probe"
	contents := "#!/bin/sh\nexit 0\n"
	if runtime.GOOS == "windows" {
		name += ".bat"
		contents = "@echo off\n"
	}
	for _, directory := range []string{first, second} {
		if err := os.WriteFile(filepath.Join(directory, name), []byte(contents), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(first, "directory"), 0o755); err != nil {
		t.Fatal(err)
	}
	oldPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", filepath.Join(first, "missing")+string(os.PathListSeparator)+first+string(os.PathListSeparator)+second); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })

	matches := getPathExecutables()
	count := 0
	for _, match := range matches {
		if match == name {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("duplicate executable count = %d in %#v", count, matches)
	}
}

func TestWordCompleterCompletionModes(t *testing.T) {
	completer := &wordCompleter{commands: []string{"cat", "car", "echo"}}

	newLine, length := completer.Do([]rune("c"), 1)
	if string(newLine[0]) != "a" || length != 1 {
		t.Fatalf("shared-prefix completion = %#v, length %d", newLine, length)
	}

	newLine, length = completer.Do([]rune("cat"), 3)
	if !reflect.DeepEqual(newLine, [][]rune{[]rune(" ")}) || length != 3 {
		t.Fatalf("single completion = %#v, length %d", newLine, length)
	}

	output := captureOutput(t, func() {
		newLine, length = completer.Do([]rune("ca"), 2)
	})
	if len(newLine) != 2 || output != "\x07" || length != 2 {
		t.Fatalf("multiple completions = %#v, output %q, length %d", newLine, output, length)
	}

	output = captureOutput(t, func() {
		newLine, length = completer.Do([]rune("echo "), 5)
	})
	if newLine != nil || output != "\x07" || length != 0 {
		t.Fatalf("argument completion = %#v, output %q, length %d", newLine, output, length)
	}
}

func TestRunBuiltinUnknownAndTypeWithoutArgument(t *testing.T) {
	if runBuiltin("unknown", nil, Redirection{}) {
		t.Fatal("unknown command reported as builtin")
	}
	if !runBuiltin("type", nil, Redirection{}) {
		t.Fatal("type without argument was not handled as builtin")
	}
}
