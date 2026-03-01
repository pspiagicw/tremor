package typechecker

import (
	"os"
	"strings"
	"testing"

	"github.com/pspiagicw/tremor/diagnostic"
	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/types"
)

func TestTypecheckerDiagnosticMessageAndLocation(t *testing.T) {
	_ = os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	input := "let a int = 1 if a then end"
	l := lexer.NewLexerWithFile(input, "typecheck.tm")
	p := parser.NewParser(l)
	ast := p.ParseAST()

	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected parser errors: %v", p.Errors())
	}

	tc := NewTypeChecker()
	tc.SetSourceContext("typecheck.tm", input)
	scope := NewScope()
	scope.SetupBuiltinFunctions()

	got := tc.TypeCheck(ast, scope)
	if got != types.UnknownType {
		t.Fatalf("expected unknown type from failed typecheck, got %s", got)
	}

	errs := tc.Errors()
	if len(errs) == 0 {
		t.Fatalf("expected typechecker errors, got none")
	}

	first := errs[0]
	if first.Error() != "If condition must be bool, got int." {
		t.Fatalf("unexpected typechecker message: %q", first.Error())
	}

	rendered := diagnostic.StripANSI(diagnostic.Render(first))
	if !strings.Contains(rendered, "--> typecheck.tm:1:18") {
		t.Fatalf("expected typechecker location in diagnostic, got: %s", rendered)
	}
	if !strings.Contains(rendered, "1 | let a int = 1 if a then end") {
		t.Fatalf("expected typechecker source line in diagnostic, got: %s", rendered)
	}
}
