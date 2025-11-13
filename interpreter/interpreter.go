package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pspiagicw/fenc/dump"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/tremor/compiler"
	"github.com/pspiagicw/tremor/lexer"
	"github.com/pspiagicw/tremor/parser"
	"github.com/pspiagicw/tremor/typechecker"
	"github.com/pspiagicw/tremor/types"
)

func StartREPL() {
	reader := bufio.NewReader(os.Stdin)
	emptyScope := typechecker.NewScope()
	for {
		fmt.Print(">>> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			goreland.LogFatal("Error reading input: %v", err)
		}
		value := strings.TrimSpace(input)

		l := lexer.NewLexer(value)
		p := parser.NewParser(l)
		ast := p.ParseAST()

		if len(p.Errors()) != 0 {
			goreland.LogError("Parser has errors!")
			for i, err := range p.Errors() {
				goreland.LogError("ERROR %d: %s", i, err)
			}
			continue
		}

		fmt.Printf("AST: %s\n", ast.String())
		t := typechecker.NewTypeChecker()

		valueType := t.TypeCheck(ast, emptyScope)

		if len(t.Errors()) != 0 {
			goreland.LogError("Typechecker has errors!")
			for i, err := range t.Errors() {
				goreland.LogError("ERROR %d: %s", i, err)
			}
			continue
		}

		if valueType == types.UnknownType {
			goreland.LogError("Typecheck failed!")
			continue
		}

		fmt.Printf("TYPE: %s\n", valueType)

		typeMap := t.Map()
		c := compiler.NewCompiler(typeMap)
		// TODO: Add a err value to the compiler!
		c.Compile(ast)
		bytecode := c.Bytecode()
		dump.Constants(bytecode.Constants)
		dump.Dump(bytecode.Tape)
	}

	os.Exit(0)
}
