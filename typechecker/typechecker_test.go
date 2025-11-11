package typechecker

import (
	"testing"

	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/types"
	"github.com/stretchr/testify/assert"
)

func TestParenthesisExpression(t *testing.T) {
	input := `(1 + 2) * (3 * 3)`

	expected := types.IntType

	testTypeChecking(t, input, expected)

}

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

func TestLetStatementWithoutType(t *testing.T) {
	input := `let a = true`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestLetStatementInt(t *testing.T) {
	input := `let a int = 1`

	expected := types.IntType

	testTypeChecking(t, input, expected)
}

func TestLambdaExpression(t *testing.T) {
	input := `let a fn() void = fn() void then print("Hello!") end`

	expected := types.NewFunctionType([]*types.Type{}, types.VoidType)

	testTypeChecking(t, input, expected)
}

func TestLambdaExpressionWithReturnType(t *testing.T) {
	input := `let a fn() int = fn() int then return 0 end`

	expected := types.NewFunctionType([]*types.Type{}, types.IntType)

	testTypeChecking(t, input, expected)
}

func TestLambdaExpressionWithArgsAndReturnType(t *testing.T) {
	input := `let a fn(int, int) int = fn(a int, b int) int then return a + b end`

	expected := types.NewFunctionType([]*types.Type{types.IntType, types.IntType}, types.IntType)

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
	input := `fn adder(x int, y int) (fn(int) int) then return fn(a int) int then return a + y end end`

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

func TestIntAddition(t *testing.T) {
	input := `1 + 2`

	expected := types.IntType

	testTypeChecking(t, input, expected)
}

func TestFloatAddition(t *testing.T) {
	input := `1.5 + 2.3`

	expected := types.FloatType

	testTypeChecking(t, input, expected)
}

func TestStringConcatenation(t *testing.T) {
	input := `"hello" .. " world"`

	expected := types.StringType

	testTypeChecking(t, input, expected)
}

func TestBooleanAnd(t *testing.T) {
	input := `true and false`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestBooleanOr(t *testing.T) {
	input := `true or false`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestIntEquality(t *testing.T) {
	input := `1 == 2`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestFloatEquality(t *testing.T) {
	input := `1.1 == 2.2`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestStringEquality(t *testing.T) {
	input := `"a" == "b"`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestStringNotEquality(t *testing.T) {
	input := `"a" != "b"`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestIntLessThan(t *testing.T) {
	input := `1 < 2`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestIntLessThanEqual(t *testing.T) {
	input := `1 <= 2`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestFloatLessThan(t *testing.T) {
	input := `1.5 < 2.5`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestFloatLessThanEqual(t *testing.T) {
	input := `1.5 <= 2.5`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestIntGreaterThan(t *testing.T) {
	input := `1 > 2`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestIntGreaterThanEqual(t *testing.T) {
	input := `1 >= 2`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestFloatGreaterThan(t *testing.T) {
	input := `1.5 > 2.5`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestFloatGreaterThanEqual(t *testing.T) {
	input := `1.5 >= 2.5`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestIntPlusFloat(t *testing.T) {
	input := `1 + 2.5`

	expected := types.FloatType

	testTypeChecking(t, input, expected)
}

func TestFloatPlusInt(t *testing.T) {
	input := `2.5 + 1`

	expected := types.FloatType

	testTypeChecking(t, input, expected)
}

func TestIntMinusFloat(t *testing.T) {
	input := `10 - 3.2`

	expected := types.FloatType

	testTypeChecking(t, input, expected)
}

func TestFloatMinusInt(t *testing.T) {
	input := `3.2 - 10`

	expected := types.FloatType

	testTypeChecking(t, input, expected)
}

func TestIntTimesFloat(t *testing.T) {
	input := `4 * 2.5`

	expected := types.FloatType

	testTypeChecking(t, input, expected)
}

func TestFloatTimesInt(t *testing.T) {
	input := `2.5 * 4`

	expected := types.FloatType

	testTypeChecking(t, input, expected)
}

func TestIntDivideFloat(t *testing.T) {
	input := `10 / 2.5`

	expected := types.FloatType

	testTypeChecking(t, input, expected)
}

func TestFloatDivideInt(t *testing.T) {
	input := `2.5 / 5`

	expected := types.FloatType

	testTypeChecking(t, input, expected)
}

func TestIntEqualsFloat(t *testing.T) {
	input := `1 == 1.0`

	expected := types.BoolType

	testTypeChecking(t, input, expected)
}

func TestFloatLessThanInt(t *testing.T) {
	input := `2.5 < 3`

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

	assert.Equal(t, expected.Kind, got.Kind, "Expected correct type.")

	if got.Kind == types.FUNCTION {
		assert.Equal(t, expected.ReturnType, got.ReturnType, "Return type don't match")
		for i := range expected.Args {
			assert.Equal(t, expected.Args[i], got.Args[i], "Args for type don't match.")
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
