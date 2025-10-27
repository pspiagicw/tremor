package compiler

import (
	"testing"

	"github.com/pspiagicw/fenc/code"
	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/typechecker"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {

	input := `1`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
	}

	testCompiler(t, input, expected)
}

func TestAdd(t *testing.T) {
	input := `1 + 2`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.ADD_INT},
	}

	testCompiler(t, input, expected)
}

// --------------------------------------------
// Integer Arithmetic
// --------------------------------------------

func TestAddInt(t *testing.T) {
	input := `1 + 2`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.ADD_INT},
	}

	testCompiler(t, input, expected)
}

func TestSubInt(t *testing.T) {
	input := `5 - 3`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.SUB_INT},
	}

	testCompiler(t, input, expected)
}

func TestMulInt(t *testing.T) {
	input := `2 * 4`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.MUL_INT},
	}

	testCompiler(t, input, expected)
}

func TestDivInt(t *testing.T) {
	input := `10 / 5`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.DIV_INT},
	}

	testCompiler(t, input, expected)
}

// --------------------------------------------
// Float Arithmetic
// --------------------------------------------

func TestAddFloat(t *testing.T) {
	input := `1.0 + 2.0`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.ADD_FLOAT},
	}

	testCompiler(t, input, expected)
}

func TestSubFloat(t *testing.T) {
	input := `5.5 - 1.5`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.SUB_FLOAT},
	}

	testCompiler(t, input, expected)
}

func TestMulFloat(t *testing.T) {
	input := `2.5 * 3.0`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.MUL_FLOAT},
	}

	testCompiler(t, input, expected)
}

func TestDivFloat(t *testing.T) {
	input := `10.0 / 2.0`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.DIV_FLOAT},
	}

	testCompiler(t, input, expected)
}

// --------------------------------------------
// Boolean Expressions
// --------------------------------------------

func TestBoolAnd(t *testing.T) {
	input := `true and false`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.AND_BOOL},
	}

	testCompiler(t, input, expected)
}

func TestBoolOr(t *testing.T) {
	input := `true or false`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.OR_BOOL},
	}

	testCompiler(t, input, expected)
}

func TestBoolNot(t *testing.T) {
	t.Skip()
	input := `not true`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		// {OpCode: code.NOT},
	}

	testCompiler(t, input, expected)
}

// --------------------------------------------
// Comparisons
// --------------------------------------------

func TestCompareInts(t *testing.T) {
	cases := []struct {
		input string
		op    code.Op
	}{
		{"1 < 2", code.LT_INT},
		{"1 <= 2", code.LTE_INT},
		{"1 > 2", code.GT_INT},
		{"1 >= 2", code.GTE_INT},
		{"1 == 2", code.EQ},
		{"1 != 2", code.NEQ},
	}

	for _, c := range cases {
		expected := []code.Instruction{
			{OpCode: code.PUSH, Args: []int{0}},
			{OpCode: code.PUSH, Args: []int{1}},
			{OpCode: c.op},
		}
		testCompiler(t, c.input, expected)
	}
}

func TestCompareFloats(t *testing.T) {
	cases := []struct {
		input string
		op    code.Op
	}{
		{"1.0 < 2.0", code.LT_FLOAT},
		{"1.0 <= 2.0", code.LTE_FLOAT},
		{"1.0 > 2.0", code.GT_FLOAT},
		{"1.0 >= 2.0", code.GTE_FLOAT},
	}

	for _, c := range cases {
		expected := []code.Instruction{
			{OpCode: code.PUSH, Args: []int{0}},
			{OpCode: code.PUSH, Args: []int{1}},
			{OpCode: c.op},
		}
		testCompiler(t, c.input, expected)
	}
}

// --------------------------------------------
// String Operations
// --------------------------------------------

func TestStringConcatDoubleQuotes(t *testing.T) {
	input := `"hello" .. "world"`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.ADD_STRING},
	}

	testCompiler(t, input, expected)
}

func TestStringConcatSingleQuotes(t *testing.T) {
	input := `'foo' .. 'bar'`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.ADD_STRING},
	}

	testCompiler(t, input, expected)
}

func TestStringConcatMultiline(t *testing.T) {
	input := `[[hello]] .. [[world]]`

	expected := []code.Instruction{
		{OpCode: code.PUSH, Args: []int{0}},
		{OpCode: code.PUSH, Args: []int{1}},
		{OpCode: code.ADD_STRING},
	}

	testCompiler(t, input, expected)
}

func testCompiler(t *testing.T, input string, expected []code.Instruction) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	tc := typechecker.NewTypeChecker()

	ast := p.ParseAST()

	assert.Empty(t, p.Errors(), "Parser has errors!")

	scope := typechecker.NewScope()
	scope.SetupBuiltinFunctions()

	_ = tc.TypeCheck(ast, scope)
	assert.Empty(t, tc.Errors(), "Type Checker has errors!")

	cmp := NewCompiler(tc, scope)
	cmp.Compile(ast)

	bytecode := cmp.Bytecode()

	assert.Equal(t, expected, bytecode.Tape, "Compiled code differs!")
}
