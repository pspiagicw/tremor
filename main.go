package main

import (
	"os"

	"github.com/pspiagicw/fenc/convert"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/compiler"
	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/typechecker"
)

func main() {
	if len(os.Args) != 2 {
		goreland.LogFatal("Expected 1 arguments program")
	}

	program := os.Args[1]

	code := readFile(program)
	AST := parseFile(code)

	c := compiler.NewCompiler()
	c.Compile(AST)

	bytecode := c.Bytecode()
	content := convert.Convert(bytecode.Tape, bytecode.Constants)
	err := os.WriteFile(program+".bc", content, os.ModeAppend)
	if err != nil {
		goreland.LogFatal("Can't write to file")
	}

}

func readFile(program string) string {
	content, err := os.ReadFile(program)

	if err != nil {
		goreland.LogFatal("Error reading file: %v", err)
	}

	return string(content)
}

func parseFile(code string) ast.Node {
	l := lexer.NewLexer(code)
	p := parser.NewParser(l)

	ast := p.ParseAST()

	if len(p.Errors()) != 0 {
		goreland.LogFatal("Parser has errors: %v", p.Errors())
	}

	tp := typechecker.NewTypeChecker()
	scope := typechecker.NewScope()
	scope.SetupBuiltinFunctions()
	_ = tp.TypeCheck(ast, scope)

	if len(tp.Errors()) != 0 {
		goreland.LogFatal("Type checker has errors: %v", p.Errors())
	}

	return ast
}
