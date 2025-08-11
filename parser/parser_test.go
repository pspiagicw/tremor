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

func testParser(t *testing.T, input string, expected string) {
	l := lexer.NewLexer(input)
	p := NewParser(l)

	node := p.ParseAST()

	result := node.String()

	if input != result {
		t.Fatalf("Expected '%s', got '%s'", expected, result)
	}
}
