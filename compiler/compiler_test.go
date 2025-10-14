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

	cmp := NewCompiler()
	cmp.Compile(ast)

	bytecode := cmp.Bytecode()

	assert.Equal(t, expected, bytecode.Tape, "Compiled code differs!")
}
