package shell

import "testing"

func TestInvalidCommandCompletionLeavesInputUnchangedAndRingsBell(t *testing.T) {
	var completions [][]rune
	var length int
	output := captureOutput(t, func() {
		completions, length = NewWordCompleter().Do([]rune("xyz"), len([]rune("xyz")))
	})

	if completions != nil {
		t.Fatalf("completions = %#v, want nil so input remains unchanged", completions)
	}
	if length != len("xyz") {
		t.Fatalf("completion length = %d, want %d", length, len("xyz"))
	}
	if output != "\x07" {
		t.Fatalf("completion output = %q, want bell character", output)
	}
}
