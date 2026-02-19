package compiler

import (
	"testing"

	"github.com/pspiagicw/fenc/vm"
	"github.com/pspiagicw/tremor/builtins"
	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/typechecker"
	"github.com/stretchr/testify/assert"
)

func TestPrintB(t *testing.T) {

	tt := []string{
		`print(1)`,
		`print("hello, world")`,
		`print(true)`,
		`print(false)`,
		`print([1,2,3,4])`,
		`print({"something": 1, "else": 2})`,
	}

	for _, testcase := range tt {
		t.Run(testcase, func(t *testing.T) {
			testBuiltin(t, testcase)
		})
	}
}

func TestLen(t *testing.T) {
	tt := []string{
		`len("something")`,
		`len([1,2,3,4])`,
		`len({"something": 2, "else": 3})`,
	}

	for _, testcase := range tt {
		t.Run(testcase, func(t *testing.T) {
			testBuiltin(t, testcase)
		})
	}
}

func testBuiltin(t *testing.T, input string) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	tc := typechecker.NewTypeChecker()

	ast := p.ParseAST()

	assert.Empty(t, p.Errors(), "Parser has errors!")

	scope := typechecker.NewScope()
	scope.SetupBuiltinFunctions()

	_ = tc.TypeCheck(ast, scope)
	result := assert.Empty(t, tc.Errors(), "Type Checker has errors!")
	if result == false {
		t.FailNow()
	}

	cmp := NewCompiler(tc.Map())
	err := cmp.Compile(ast)
	assert.Nil(t, err, "Compiler has a error!")

	bytecode := cmp.Bytecode()
	builtins := builtins.GetBuiltins()
	vm := vm.NewVM(bytecode, builtins)
	vm.Run()
}
