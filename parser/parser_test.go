package parser

import (
	"testing"

	"github.com/pspiagicw/fener/ast"
	"github.com/pspiagicw/fener/lexer"
)

func TestLetStatement(t *testing.T) {
	input := "let a = 1"

	expected := "let a = 1"

	testParser(t, input, expected)
}

func TestMultipleStatements(t *testing.T) {
	input := "let a = 1 let b = 2"

	expected := "let a = 1 let b = 2"

	testParser(t, input, expected)
}

func testParser(t *testing.T, input string, expected string) {
	l := lexer.NewLexer(input)
	p := NewParser(l)

	node := p.ParseAST()

	result := node.String()

	if input != result {
		t.Logf("Expected '%q', got '%q' ", expected, result)
		for i := range expected {
			if expected[i] != result[i] {
				t.Fatalf("Diff at %d: %q vs %q\n", i, expected[i], result[i])
				break
			}
		}

	}

}

func TestParser_parseFunctionCallExpression(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		l *lexer.Lexer
		// Named input parameters for target function.
		left ast.Expression
		want ast.Expression
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.l)
			got := p.parseFunctionCallExpression(tt.left)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("parseFunctionCallExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
