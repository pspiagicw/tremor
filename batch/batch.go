package batch

import (
	"fmt"
	"log"
	"os"

	"github.com/pspiagicw/fenc/dump"
	"github.com/pspiagicw/fenc/vm"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/compiler"
	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/typechecker"
)

func ExecFile(filename string) {
	code := readFile(filename)
	AST, typeMap := parseFile(code)

	c := compiler.NewCompiler(typeMap)
	err := c.Compile(AST)
	if err != nil {
		log.Fatalf("ERROR: %s", err)

	}

	bytecode := c.Bytecode()

	dump.Constants(bytecode.Constants)
	dump.Dump(bytecode.Tape)

	vm := vm.NewVM(bytecode)
	vm.Run()
}
func readFile(program string) string {
	content, err := os.ReadFile(program)

	if err != nil {
		goreland.LogFatal("Error reading file: %v", err)
	}

	return string(content)
}

func parseFile(code string) (ast.Node, typechecker.TypeMap) {
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
		goreland.LogError("Type checker has errors")
		for _, err := range tp.Errors() {
			fmt.Printf("ERROR: %s\n", err)
		}
		goreland.LogFatal("")
	}

	return ast, tp.Map()
}
