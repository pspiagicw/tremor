package batch

import (
	"log"
	"os"

	"github.com/pspiagicw/fenc/dump"
	"github.com/pspiagicw/fenc/vm"
	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/builtins"
	"github.com/pspiagicw/tremor/compiler"
	"github.com/pspiagicw/tremor/diagnostic"
	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/typechecker"
)

func ExecFile(filename string) {
	code := readFile(filename)
	AST, typeMap := parseFile(code, filename)

	c := compiler.NewCompiler(typeMap)
	c.SetSourceContext(filename, code)
	err := c.Compile(AST)
	if err != nil {
		log.Fatalf("%s", diagnostic.Render(err))

	}

	bytecode := c.Bytecode()

	dump.Constants(bytecode.Constants)
	dump.Dump(bytecode.Tape)

	vm := vm.NewVM(bytecode, builtins.GetBuiltins())
	vm.Run()
}
func readFile(program string) string {
	content, err := os.ReadFile(program)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	return string(content)
}

func parseFile(code string, filename string) (ast.Node, typechecker.TypeMap) {
	l := lexer.NewLexerWithFile(code, filename)
	p := parser.NewParser(l)

	ast := p.ParseAST()

	if len(p.Errors()) != 0 {
		log.Println("Parser has errors:")
		for _, err := range p.Errors() {
			log.Println(diagnostic.Render(err))
		}
	}

	tp := typechecker.NewTypeChecker()
	tp.SetSourceContext(filename, code)
	scope := typechecker.NewScope()
	scope.SetupBuiltinFunctions()
	_ = tp.TypeCheck(ast, scope)

	if len(tp.Errors()) != 0 {
		log.Println("Type checker has errors")
		for _, err := range tp.Errors() {
			log.Println(diagnostic.Render(err))
		}
		log.Fatal()
	}

	return ast, tp.Map()
}
