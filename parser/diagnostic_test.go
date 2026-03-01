package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/pspiagicw/tremor/diagnostic"
	"github.com/pspiagicw/tremor/lexer"
)

func TestParserDiagnosticMessageAndLocation(t *testing.T) {
	_ = os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	input := "let a int 1"
	l := lexer.NewLexerWithFile(input, "parser.tm")
	p := NewParser(l)

	_ = p.ParseAST()

	errs := p.Errors()
	if len(errs) == 0 {
		t.Fatalf("expected parser errors, got none")
	}

	first := errs[0]
	if first.Error() != "Expected token type ASSIGN, got INTEGER." {
		t.Fatalf("unexpected parser message: %q", first.Error())
	}

	rendered := diagnostic.StripANSI(diagnostic.Render(first))
	if !strings.Contains(rendered, "--> parser.tm:1:11") {
		t.Fatalf("expected parser location in diagnostic, got: %s", rendered)
	}
	if !strings.Contains(rendered, "1 | let a int 1") {
		t.Fatalf("expected parser source line in diagnostic, got: %s", rendered)
	}
}
