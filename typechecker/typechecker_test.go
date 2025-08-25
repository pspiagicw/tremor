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

	input := `fn hello() void then print("Hello, World") end`

	expected := types.NewFunctionType([]*types.Type{}, types.VoidType)

	testTypeChecking(t, input, expected)
}

func TestFunctionStatementWithReturnType(t *testing.T) {
	input := `fn hello() string then return "hello" end`

	expected := types.NewFunctionType([]*types.Type{}, types.StringType)

	testTypeChecking(t, input, expected)
}

func TestFunctionStatementWithArgTypes(t *testing.T) {
	input := `fn add(a int, b int) int then return a + b end`

	expected := types.NewFunctionType([]*types.Type{types.IntType, types.IntType}, types.IntType)

	testTypeChecking(t, input, expected)
}

func TestFunctionStatementWithFunctionArgTypes(t *testing.T) {
	input := `fn apply(val int, somefunc fn(int) int) int then return somefunc(val) end`

	expected := types.NewFunctionType(
		[]*types.Type{
			types.IntType,
			types.NewFunctionType(
				[]*types.Type{
					types.IntType,
				},
				types.IntType,
			),
		},
		types.IntType,
	)

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

	if got.Kind != expected.Kind {
		t.Fatalf("Expected type of %s, got type of %s", expected.Kind, got.Kind)
	}

	if got.Kind == types.FUNCTION {
		if got.ReturnType != expected.ReturnType {
			t.Fatalf("Expected return type of %s, got type of %s for function statement", expected.ReturnType.Kind, got.ReturnType.Kind)
		}
		for i := range expected.Args {
			if got.Args[i].Kind != expected.Args[i].Kind {
				t.Fatalf("Expected arg type of %s, got arg type of %s.", expected.Args[i].Kind, got.Args[i].Kind)
			}
		}
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
