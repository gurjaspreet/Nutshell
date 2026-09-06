package shell

import (
	"reflect"
	"testing"
)

func TestLexerTokenize(t *testing.T) {
	lexer := Lexer{}
	got := lexer.Tokenize(`echo "hello world" 'from shell'`)
	want := []string{"echo", "hello world", "from shell"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokenize() = %#v, want %#v", got, want)
	}
}

func TestParseRedirection(t *testing.T) {
	gotArgs, gotRedirection := ParseRedirection([]string{"echo", "hello", "2>>", "errors.log"})
	wantArgs := []string{"echo", "hello"}
	wantRedirection := Redirection{FileDescriptor: 2, Filename: "errors.log", Append: true}

	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("ParseRedirection() args = %#v, want %#v", gotArgs, wantArgs)
	}
	if gotRedirection != wantRedirection {
		t.Fatalf("ParseRedirection() redirection = %#v, want %#v", gotRedirection, wantRedirection)
	}
}
