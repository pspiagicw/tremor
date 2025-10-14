package compiler

import (
	"strconv"

	"github.com/pspiagicw/fenc/emitter"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/tremor/ast"
)

type Compiler struct {
	e *emitter.Emitter
}

func NewCompiler() *Compiler {
	return &Compiler{
		e: emitter.NewEmitter(),
	}
}

func (c *Compiler) Compile(node ast.Node) {
	switch node := node.(type) {
	case *ast.AST:
		c.compileAST(node)
	case ast.IntegerExpression:
		c.compileInteger(node)
	case ast.ExpressionStatement:
		c.Compile(node.Inside)
	case ast.BinaryExpression:
		c.compileBinary(node)
	default:
		goreland.LogFatal("Can't compile type '%v'", node)
	}
}
func (c *Compiler) compileBinary(node ast.BinaryExpression) {
}
func (c *Compiler) compileInteger(node ast.IntegerExpression) {
	value, err := strconv.Atoi(node.Value)
	if err != nil {
		goreland.LogFatal("Error converting '%s' to integer", node.Value)
	}
	c.e.PushInt(value)
}
func (c *Compiler) compileAST(node *ast.AST) {
	for _, statement := range node.Statements {
		c.Compile(statement)
	}
}
func (c *Compiler) Bytecode() emitter.ByteCode {
	return c.e.Bytecode()
}
