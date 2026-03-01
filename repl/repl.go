package repl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pspiagicw/fenc/dump"
	"github.com/pspiagicw/fenc/vm"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/tremor/builtins"
	"github.com/pspiagicw/tremor/compiler"
	"github.com/pspiagicw/tremor/diagnostic"
	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/typechecker"
	"github.com/pspiagicw/tremor/types"
)

func StartREPL() {
	debugMode := os.Getenv("TREMOR_DEBUG") == "1"

	emptyScope := typechecker.NewScope()
	emptyScope.SetupBuiltinFunctions()
	t := typechecker.NewTypeChecker()
	typeMap := t.Map()
	c := compiler.NewCompiler(typeMap)

	for {

		value := getLine()

		l := lexer.NewLexerWithFile(value, "<repl>")
		p := parser.NewParser(l)
		ast := p.ParseAST()

		if len(p.Errors()) != 0 {
			goreland.LogError("Parser has errors!")
			for _, err := range p.Errors() {
				log.Println(diagnostic.Render(err))
			}
			continue
		}

		if debugMode {
			fmt.Printf("AST: %s\n", ast.String())
		}

		t.SetSourceContext("<repl>", value)
		valueType := t.TypeCheck(ast, emptyScope)

		if len(t.Errors()) != 0 {
			goreland.LogError("Typechecker has errors!")
			for _, err := range t.Errors() {
				log.Println(diagnostic.Render(err))
			}
			// Reset typechecker messages.
			t.Flush()
			continue
		}

		if valueType == types.UnknownType {
			log.Println("Typecheck failed!")
			continue
		}

		if debugMode {
			fmt.Printf("TYPE: %s\n", valueType)
		}

		tm := t.Map()
		c.SetTypeMap(tm)
		c.SetSourceContext("<repl>", value)
		err := c.Compile(ast)

		if err != nil {
			log.Printf("Compiler faced errors:\n%s", diagnostic.Render(err))
			continue
		}

		bytecode := c.Bytecode()
		if debugMode {
			// dump.Constants(bytecode.Constants)
			dump.Dump(bytecode.Tape)
		}

		// fmt.Println("==== OUTPUT === ")

		vm := vm.NewVM(bytecode, builtins.GetBuiltins())

		vm.Run()

		fmt.Println(vm.Peek())

		// fmt.Println("---- ...... --- ")
	}
}

func getLine() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">>> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	value := strings.TrimSpace(input)

	return value
}
