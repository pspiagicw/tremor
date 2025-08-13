package parser

import (
	"testing"

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

func TestIfStatement(t *testing.T) {
	input := `if true then print("true") end`

	expected := `if true then print("true") end`

	testParser(t, input, expected)
}

func TestIfElseStatement(t *testing.T) {
	input := `if true then print("true") else print("false") end`

	expected := `if true then print("true") else print("false") end`

	testParser(t, input, expected)
}

func TestReturnStatement(t *testing.T) {
	input := `return 1`

	testParser(t, input, input)
}

func TestExpressionStatement(t *testing.T) {
	input := `1`

	expected := `1`

	testParser(t, input, expected)
}

func TestExpressionStatementComplex(t *testing.T) {
	input := `1 + 2 * 3`

	expected := `(1 + (2 * 3))`

	testParser(t, input, expected)
}
func TestFunctionStatement(t *testing.T) {
	input := `fn hello() then print("Hello, World") end`

	testParser(t, input, input)
}

func TestFunctionStatementComplex(t *testing.T) {
	input := `fn hello(a, b) then return a + b end`

	expected := `fn hello(a, b) then return (a + b) end`

	testParser(t, input, expected)

}

func testParser(t *testing.T, input string, expected string) {
	l := lexer.NewLexer(input)
	p := NewParser(l)

	node := p.ParseAST()

	result := node.String()

	if len(expected) != len(result) {
		t.Errorf("The length doesn't match, expected: %d, got: %d", len(expected), len(result))
	}

	if expected != result {
		t.Fatalf("Expected '%q', got '%q' ", expected, result)
	}

}
