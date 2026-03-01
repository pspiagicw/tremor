package diagnostic

import (
	"os"
	"strings"
	"testing"

	"github.com/pspiagicw/tremor/token"
)

func TestPrettyWithSourceSpan(t *testing.T) {
	_ = os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	src := "let x int 1\n"
	tok := &token.Token{
		Type:   token.INTEGER,
		Value:  "1",
		Offset: 10,
		Line:   1,
		Column: 11,
	}

	d := NewAtToken("parser", "sample.tm", src, tok, len(tok.Value), "Expected token type ASSIGN, got INTEGER.")
	rendered := StripANSI(d.Pretty())

	if !strings.Contains(rendered, "error[parser]: Expected token type ASSIGN, got INTEGER.") {
		t.Fatalf("missing header in rendered diagnostic: %s", rendered)
	}
	if !strings.Contains(rendered, "--> sample.tm:1:11") {
		t.Fatalf("missing source location in rendered diagnostic: %s", rendered)
	}
	if !strings.Contains(rendered, "1 | let x int 1") {
		t.Fatalf("missing source line in rendered diagnostic: %s", rendered)
	}
	if !strings.Contains(rendered, "^") {
		t.Fatalf("missing underline in rendered diagnostic: %s", rendered)
	}
}

func TestRenderFallback(t *testing.T) {
	_ = os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	d := New("typechecker", "", "", "Type mismatch.")
	rendered := StripANSI(Render(d))
	if !strings.Contains(rendered, "error[typechecker]: Type mismatch.") {
		t.Fatalf("unexpected render output: %s", rendered)
	}
}

func TestPrettyColorEnabled(t *testing.T) {
	_ = os.Unsetenv("NO_COLOR")
	_ = os.Setenv("TREMOR_COLOR", "always")
	defer os.Unsetenv("TREMOR_COLOR")

	d := New("parser", "", "", "Unexpected token.")
	rendered := d.Pretty()

	if !strings.Contains(rendered, "\x1b[") {
		t.Fatalf("expected ansi color escape sequence, got: %q", rendered)
	}
}
