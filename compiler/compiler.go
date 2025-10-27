package compiler

import (
	"strconv"

	"github.com/pspiagicw/fenc/emitter"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/token"
	"github.com/pspiagicw/tremor/typechecker"
	"github.com/pspiagicw/tremor/types"
)

type Compiler struct {
	e     *emitter.Emitter
	tp    *typechecker.TypeChecker
	scope *typechecker.TypeScope
}

func NewCompiler(tp *typechecker.TypeChecker, scope *typechecker.TypeScope) *Compiler {
	return &Compiler{
		e:     emitter.NewEmitter(),
		tp:    tp,
		scope: scope,
	}
}

func (c *Compiler) Compile(node ast.Node) {
	switch node := node.(type) {
	case *ast.AST:
		c.compileAST(node)
	case ast.IntegerExpression:
		c.compileInteger(node)
	case ast.FloatExpression:
		c.compileFloat(node)
	case ast.ExpressionStatement:
		c.Compile(node.Inside)
	case ast.BinaryExpression:
		c.compileBinary(node)
	default:
		goreland.LogFatal("Can't compile type '%v'", node)
	}
}
func (c *Compiler) compileFloat(node ast.FloatExpression) {
	value, err := strconv.ParseFloat(node.Value, 32)
	if err != nil {
		goreland.LogFatal("Error converting '%s' to float", node.Value)
	}
	c.e.PushFloat(float32(value))
}
func (c *Compiler) compileBinary(node ast.BinaryExpression) {
	operator := node.Operator.Type
	nodeType := c.tp.TypeCheck(node, c.scope)

	c.Compile(node.Left)
	c.Compile(node.Right)

	switch operator {
	case token.PLUS:
		c.emitPlus(nodeType)
	}

}
func (c *Compiler) emitPlus(nodeType *types.Type) {
	switch nodeType {
	case types.IntType:
		c.e.AddInt()
	case types.FloatType:
		c.e.AddFloat()
	}
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
