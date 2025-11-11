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
	e       *emitter.Emitter
	typeMap typechecker.TypeMap
}

func NewCompiler(typeMap typechecker.TypeMap) *Compiler {
	return &Compiler{
		e:       emitter.NewEmitter(),
		typeMap: typeMap,
	}
}

func (c *Compiler) Compile(node ast.Node) {
	switch node := node.(type) {
	case *ast.AST:
		c.compileAST(node)
	case *ast.IntegerExpression:
		c.compileInteger(node)
	case *ast.FloatExpression:
		c.compileFloat(node)
	case *ast.ExpressionStatement:
		c.Compile(node.Inside)
	case *ast.BinaryExpression:
		c.compileBinary(node)
	case *ast.BooleanExpression:
		c.compileBoolean(node)
	case *ast.StringExpression:
		c.compileString(node)
	case *ast.ParenthesisExpression:
		c.compileParenthesis(node)
	default:
		goreland.LogFatal("Can't compile type '%v'", node)
	}
}
func (c *Compiler) compileParenthesis(node *ast.ParenthesisExpression) {
	c.Compile(node.Inside)
}
func (c *Compiler) compileString(node *ast.StringExpression) {
	value := node.Value

	c.e.PushString(value)
}
func (c *Compiler) compileBoolean(node *ast.BooleanExpression) {
	value := false
	if node.Value.Value == "true" {
		value = true
	}

	c.e.PushBool(value)
}
func (c *Compiler) compileFloat(node *ast.FloatExpression) {
	value, err := strconv.ParseFloat(node.Value, 32)
	if err != nil {
		goreland.LogFatal("Error converting '%s' to float", node.Value)
	}
	c.e.PushFloat(float32(value))
}

func (c *Compiler) compileArithmetic(node *ast.BinaryExpression) {
	returnType := c.typeMap[node]
	leftType := c.typeMap[node.Left]
	rightType := c.typeMap[node.Right]

	c.Compile(node.Left)
	if leftType == types.IntType && returnType == types.FloatType {
		c.e.ToFloat()
	}
	c.Compile(node.Right)
	if rightType == types.IntType && returnType == types.FloatType {
		c.e.ToFloat()
	}

	switch node.Operator.Type {
	case token.PLUS:
		c.emitPlus(returnType)
	case token.MINUS:
		c.emitMinus(returnType)
	case token.MULTIPLY:
		c.emitMultiply(returnType)
	case token.SLASH:
		c.emitSlash(returnType)
	}
}
func (c *Compiler) compileComparison(node *ast.BinaryExpression) {}

func (c *Compiler) compileLogical(node *ast.BinaryExpression) {
	operator := node.Operator.Type

	c.Compile(node.Left)
	c.Compile(node.Right)

	switch operator {
	case token.AND:
		c.e.AndBool()
	case token.OR:
		c.e.OrBool()
	}
}

func (c *Compiler) compileBinary(node *ast.BinaryExpression) {
	operator := node.Operator.Type

	switch operator {
	case token.PLUS, token.MINUS, token.MULTIPLY, token.SLASH:
		c.compileArithmetic(node)
	case token.GT, token.GTE, token.LT, token.LTE:
		c.compileComparison(node)
	case token.AND, token.OR:
		c.compileLogical(node)
	case token.CONCAT:
		// TODO: Maybe expand into separate function.
		c.Compile(node.Left)
		c.Compile(node.Right)
		c.e.AddString()
	default:
		goreland.LogFatal("Can't compile binary operator: %s", operator)
	}

	// TODO: This works for arithmetic, for comparison, the returntype is always boolean, thus can't emit proper instructions.
	// returnType := c.typeMap[node]
	// leftType := c.typeMap[node.Left]
	// rightType := c.typeMap[node.Right]
	//
	// // TODO: If it's a arithmetic , see if int needs to be converted to float or something.
	// // And in comparison the returnType will be boolean, you will have to compare manually.
	// c.Compile(node.Left)
	// if leftType == types.IntType && returnType == types.FloatType {
	// 	c.e.ToFloat()
	// }
	// c.Compile(node.Right)
	// if rightType == types.IntType && returnType == types.FloatType {
	// 	c.e.ToFloat()
	// }

	// switch operator {
	// case token.PLUS:
	// 	c.emitPlus(returnType)
	// case token.MINUS:
	// 	c.emitMinus(returnType)
	// case token.MULTIPLY:
	// 	c.emitMultiply(returnType)
	// case token.SLASH:
	// 	c.emitSlash(returnType)
	// case token.CONCAT:
	// 	c.e.AddString()
	// case token.OR:
	// 	c.e.OrBool()
	// case token.AND:
	// 	c.e.AndBool()
	// case token.LT:
	// 	c.Lt(node)
	// }
}
func (c *Compiler) Lt(node *ast.BinaryExpression) {
}
func (c *Compiler) emitSlash(nodeType *types.Type) {
	switch nodeType {
	case types.IntType:
		c.e.DivInt()
	case types.FloatType:
		c.e.DivFloat()
	}
}
func (c *Compiler) emitMultiply(nodeType *types.Type) {
	switch nodeType {
	case types.IntType:
		c.e.MulInt()
	case types.FloatType:
		c.e.MulFloat()
	}
}
func (c *Compiler) emitMinus(nodeType *types.Type) {
	switch nodeType {
	case types.IntType:
		c.e.SubInt()
	case types.FloatType:
		c.e.SubFloat()
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
func (c *Compiler) compileInteger(node *ast.IntegerExpression) {
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
