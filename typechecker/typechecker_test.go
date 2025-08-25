package parser

import (
	"testing"

	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/types"
)

func TestIntType(t *testing.T) {
	input := `1`

	expected := types.IntType

	testTypeChecking(t, input, expected)
}

func TestStringType(t *testing.T) {
	input := `"hello"`

	expected := types.StringType

	testTypeChecking(t, input, expected)
}

func TestBooleanType(t *testing.T) {
	input := `true`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestLetStatementBool(t *testing.T) {
	input := `let a bool = true`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestLetStatementInt(t *testing.T) {
	input := `let a = 1`

	expected := types.IntType

	testTypeChecking(t, input, expected)
}

func TestFunctionStatement(t *testing.T) {

	input := `fn hello() then print("Hello, World") end`

	expected := types.NewFunctionType([]types.Type{}, types.VoidType)

	testTypeChecking(t, input, expected)
}

func TestLetStatementString(t *testing.T) {
	input := `let b = "name"`

	expected := types.StringType

	testTypeChecking(t, input, expected)
}
func TestReturnStatement(t *testing.T) {
	input := `return 1`

	expected := types.IntType

	testTypeChecking(t, input, expected)
}
func TestIfStatement(t *testing.T) {
	input := `let a = true if a then end`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func testTypeChecking(t *testing.T, input string, expected *types.Type) {

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	typechecker := NewTypeChcker()

	ast := p.ParseAST()

	printParserErrors(t, p)

	got := typechecker.TypeCheck(ast)

	printTypeCheckerErrors(t, typechecker)

	if got != expected {
		t.Fatalf("Expected type of %s, got type of %s", expected.Kind, got.Kind)
	}
}
func printTypeCheckerErrors(t *testing.T, typechecker *TypeChecker) {
	errs := typechecker.Errors()

	if len(errs) != 0 {
		t.Errorf("The typechecker had %d errors", len(errs))
		for _, error := range errs {
			t.Errorf("%q", error)
		}
		t.Fatal()
	}
}
func printParserErrors(t *testing.T, p *parser.Parser) {

	errs := p.Errors()

	if len(errs) != 0 {
		t.Errorf("The parser had %d errors", len(errs))
		for _, error := range errs {
			t.Errorf("%q", error)
		}
		t.Fatal()
	}
}
