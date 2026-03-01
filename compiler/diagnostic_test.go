package compiler

import (
	"os"
	"strings"
	"testing"

	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/diagnostic"
	"github.com/pspiagicw/tremor/token"
	"github.com/pspiagicw/tremor/typechecker"
)

func TestCompilerDiagnosticForUnsupportedOperator(t *testing.T) {
	_ = os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	source := "1 % 2"
	node := &ast.BinaryExpression{
		Left:  &ast.IntegerExpression{Value: "1"},
		Right: &ast.IntegerExpression{Value: "2"},
		Operator: &token.Token{
			Type:   token.MODULUS,
			Value:  "%",
			Offset: 2,
			Line:   1,
			Column: 3,
		},
	}

	c := NewCompiler(typechecker.TypeMap{})
	c.SetSourceContext("compile.tm", source)

	err := c.Compile(node)
	if err == nil {
		t.Fatalf("expected compiler error, got nil")
	}
	if err.Error() != "Cannot compile binary operator MODULUS." {
		t.Fatalf("unexpected compiler message: %q", err.Error())
	}

	rendered := diagnostic.StripANSI(diagnostic.Render(err))
	if !strings.Contains(rendered, "--> compile.tm:1:3") {
		t.Fatalf("expected compiler location in diagnostic, got: %s", rendered)
	}
	if !strings.Contains(rendered, "1 | 1 % 2") {
		t.Fatalf("expected compiler source line in diagnostic, got: %s", rendered)
	}
}
