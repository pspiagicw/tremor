package typechecker

import (
	"testing"

	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/types"
	"github.com/stretchr/testify/assert"
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
	input := `let a int = 1`

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
func TestFunctionStatementWithFunctionReturnTypes(t *testing.T) {
	// TODO: Add support for lambdas to cover this.
	input := `fn adder(x int, y int) (fn(int) int) then return "something" end`

	expected := types.NewFunctionType(
		[]*types.Type{
			types.IntType,
			types.IntType,
		},
		types.NewFunctionType(
			[]*types.Type{types.IntType},
			types.IntType,
		),
	)

	testTypeChecking(t, input, expected)
}

func TestLetStatementString(t *testing.T) {
	input := `let b string = "name"`

	expected := types.StringType

	testTypeChecking(t, input, expected)
}
func TestReturnStatement(t *testing.T) {
	input := `return 1`

	expected := &types.Type{Kind: types.RETURN}

	testTypeChecking(t, input, expected)
}
func TestIfStatement(t *testing.T) {
	input := `let a bool = true if a then end`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func testTypeChecking(t *testing.T, input string, expected *types.Type) {

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	typechecker := NewTypeChecker()

	ast := p.ParseAST()

	printParserErrors(t, p)

	scope := NewScope()
	scope.SetupBuiltinFunctions()

	got := typechecker.TypeCheck(ast, scope)

	printTypeCheckerErrors(t, typechecker)

	assert.Equal(t, got.Kind, expected.Kind, "Expected correct type.")

	if got.Kind == types.FUNCTION {
		assert.Equal(t, got.ReturnType, expected.ReturnType, "Return type don't match")
		for i := range expected.Args {
			assert.Equal(t, got.Args[i], expected.Args[i], "Args for type don't match.")
		}
	}
}
func printTypeCheckerErrors(t *testing.T, typechecker *TypeChecker) {
	errs := typechecker.Errors()

	assert.Empty(t, errs, "Typechecker has errors!")
}
func printParserErrors(t *testing.T, p *parser.Parser) {

	errs := p.Errors()

	assert.Empty(t, errs, "Parser has errors!")
}
