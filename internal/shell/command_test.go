package shell

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestInvalidCommandsPrintCommandNotFound(t *testing.T) {
	got := captureOutput(t, func() {
		shellApp := New()
		for _, command := range []string{"invalid_command_1", "invalid_command_2", "invalid_command_3"} {
			if shellApp.ExecuteLine(command) {
				t.Fatalf("ExecuteLine(%q) requested exit", command)
			}
		}
	})

	want := "invalid_command_1: command not found\n" +
		"invalid_command_2: command not found\n" +
		"invalid_command_3: command not found\n"
	if got != want {
		t.Fatalf("invalid command output = %q, want %q", got, want)
	}
}

func TestExitTerminatesShell(t *testing.T) {
	if !New().ExecuteLine("exit") {
		t.Fatal("ExecuteLine(\"exit\") = false, want true")
	}
}

func TestEchoPrintsAllArguments(t *testing.T) {
	got := captureOutput(t, func() {
		New().ExecuteLine("echo hello world")
		New().ExecuteLine("echo pineapple strawberry")
	})

	want := "hello world\npineapple strawberry\n"
	if got != want {
		t.Fatalf("echo output = %q, want %q", got, want)
	}
}

func TestTypeIdentifiesBuiltinsAndUnknownCommands(t *testing.T) {
	got := captureOutput(t, func() {
		shellApp := New()
		shellApp.ExecuteLine("type echo")
		shellApp.ExecuteLine("type exit")
		shellApp.ExecuteLine("type type")
		shellApp.ExecuteLine("type invalid_command")
	})

	want := "echo is a shell builtin\n" +
		"exit is a shell builtin\n" +
		"type is a shell builtin\n" +
		"invalid_command: not found\n"
	if got != want {
		t.Fatalf("type output = %q, want %q", got, want)
	}
}

func TestTypeFindsExecutableInPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("execute permission is not represented by Unix file modes on Windows")
	}

	commandPath := createTestCommand(t, "path_command", "#!/bin/sh\nexit 0\n", 0o755)
	withPath(t, filepath.Dir(commandPath), func() {
		got := captureOutput(t, func() {
			New().ExecuteLine("type path_command")
		})

		want := "path_command is " + commandPath + "\n"
		if got != want {
			t.Fatalf("type PATH output = %q, want %q", got, want)
		}
	})
}

func TestTypeSkipsNonExecutableFiles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("execute permission is not represented by Unix file modes on Windows")
	}

	commandPath := createTestCommand(t, "not_executable", "#!/bin/sh\nexit 0\n", 0o644)
	withPath(t, filepath.Dir(commandPath), func() {
		got := captureOutput(t, func() {
			New().ExecuteLine("type not_executable")
		})

		if got != "not_executable: not found\n" {
			t.Fatalf("type non-executable output = %q", got)
		}
	})
}

func TestExternalCommandReceivesArguments(t *testing.T) {
	commandName := "custom_exe_1234"
	commandContents := "#!/bin/sh\nprintf 'Program was passed %s args (including program name).\\n' \"$#\"\nprintf 'Arg #0 (program name): %s\\n' \"$0\"\nprintf 'Arg #1: %s\\n' \"$1\"\nprintf 'Program Signature: 5998595441\\n'\n"
	if runtime.GOOS == "windows" {
		commandContents = "@echo off\necho Program was passed 2 args (including program name).\necho Arg #0 (program name): %~f0\necho Arg #1: %1\necho Program Signature: 5998595441\n"
	}

	commandPath := createTestCommand(t, commandName, commandContents, 0o755)
	withPath(t, filepath.Dir(commandPath), func() {
		got := captureOutput(t, func() {
			New().ExecuteLine(commandName + " alice")
		})
		got = strings.ReplaceAll(got, "\r\n", "\n")

		if !strings.Contains(got, "Program was passed 2 args (including program name).\n") ||
			!strings.Contains(got, "Arg #1: alice\n") ||
			!strings.Contains(got, "Program Signature: 5998595441\n") {
			t.Fatalf("external command output = %q", got)
		}
	})
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	originalStdout := os.Stdout
	originalStderr := os.Stderr
	os.Stdout = writer
	os.Stderr = writer
	defer func() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	fn()
	writer.Close()
	output, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		t.Fatal(err)
	}
	return string(output)
}

func createTestCommand(t *testing.T, name, contents string, mode os.FileMode) string {
	t.Helper()
	directory := t.TempDir()
	filename := name
	if runtime.GOOS == "windows" {
		filename += ".bat"
	}
	path := filepath.Join(directory, filename)
	if err := os.WriteFile(path, []byte(contents), mode); err != nil {
		t.Fatal(err)
	}
	return path
}

func withPath(t *testing.T, directory string, fn func()) {
	t.Helper()
	oldPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", directory+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatal(err)
	}
	defer os.Setenv("PATH", oldPath)
	fn()
}
