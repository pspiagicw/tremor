package compiler

import (
	"fmt"
	"strconv"

	"github.com/pspiagicw/fenc/emitter"
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
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.AST:
		return c.compileAST(node)
	case *ast.IntegerExpression:
		return c.compileInteger(node)
	case *ast.FloatExpression:
		return c.compileFloat(node)
	case *ast.ExpressionStatement:
		return c.Compile(node.Inside)
	case *ast.BinaryExpression:
		return c.compileBinary(node)
	case *ast.BooleanExpression:
		return c.compileBoolean(node)
	case *ast.StringExpression:
		return c.compileString(node)
	case *ast.ParenthesisExpression:
		return c.compileParenthesis(node)
	case *ast.LetStatement:
		return c.compileLetStatement(node)
	case *ast.IfStatement:
		return c.compileIfStatement(node)
	case *ast.BlockStatement:
		return c.compileBlockStatement(node)
	case *ast.IdentifierExpression:
		return c.compileIdentifierExpression(node)
	case *ast.ReturnStatement:
		return c.compileReturnStatement(node)
	case *ast.FunctionStatement:
		return c.compileFunctionStatement(node)
	case *ast.FunctionCallExpression:
		return c.compileFunctionCall(node)
	case *ast.LambdaExpression:
		return c.compileLambdaExpression(node)
	default:
		return fmt.Errorf("Can't compile type: %v", node)
	}
}
func (c *Compiler) compileLambdaExpression(node *ast.LambdaExpression) error {
	args := []string{}

	for _, arg := range node.Args {
		args = append(args, arg.Value)
	}

	return c.e.Lambda(args, func(e *emitter.Emitter) error {
		oldEmitter := c.e

		c.e = e
		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		c.e = oldEmitter

		return nil
	})
}
func (c *Compiler) compileFunctionCall(node *ast.FunctionCallExpression) error {

	for _, arg := range node.Arguments {
		err := c.Compile(arg)
		if err != nil {
			return err
		}
	}

	// TODO: Right now function name is just a expression.
	c.e.Load(node.Caller.String())

	c.e.Call(len(node.Arguments))

	return nil
}

func (c *Compiler) compileFunctionStatement(node *ast.FunctionStatement) error {
	args := []string{}
	for _, arg := range node.Args {
		args = append(args, arg.Value)
	}
	return c.e.Function(node.Name.Value, args, func(e *emitter.Emitter) error {
		oldEmitter := c.e
		// New emitter for the function's sake
		c.e = e
		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		c.e = oldEmitter

		return nil
	})
}
func (c *Compiler) compileReturnStatement(node *ast.ReturnStatement) error {
	err := c.Compile(node.Value)

	if err != nil {
		return err
	}

	c.e.ReturnValue()

	return nil
}
func (c *Compiler) compileIdentifierExpression(node *ast.IdentifierExpression) error {
	c.e.Load(node.Value.Value)

	return nil
}
func (c *Compiler) compileBlockStatement(node *ast.BlockStatement) error {
	for _, statement := range node.Statements {
		err := c.Compile(statement)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *Compiler) compileIfStatement(node *ast.IfStatement) error {
	// TODO: Change the emitter to return errors
	return c.e.If(
		func(e *emitter.Emitter) error {
			err := c.Compile(node.Condition)
			if err != nil {
				return err
			}
			return nil
		},
		func(e *emitter.Emitter) error {
			err := c.Compile(node.Consequence)
			if err != nil {
				return err
			}
			return nil
		},
		func(e *emitter.Emitter) error {
			if node.Alternative != nil {
				err := c.Compile(node.Alternative)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)

}
func (c *Compiler) compileLetStatement(node *ast.LetStatement) error {
	err := c.Compile(node.Value)
	if err != nil {
		return err
	}

	c.e.Store(node.Name.Value)
	return nil

}
func (c *Compiler) compileParenthesis(node *ast.ParenthesisExpression) error {
	return c.Compile(node.Inside)
}
func (c *Compiler) compileString(node *ast.StringExpression) error {
	value := node.Value

	c.e.PushString(value)
	return nil
}
func (c *Compiler) compileBoolean(node *ast.BooleanExpression) error {
	value := false
	if node.Value.Value == "true" {
		value = true
	}

	c.e.PushBool(value)
	return nil
}
func (c *Compiler) compileFloat(node *ast.FloatExpression) error {
	value, err := strconv.ParseFloat(node.Value, 32)
	if err != nil {
		return fmt.Errorf("Error converting '%s' to float", node.Value)
	}
	c.e.PushFloat(float32(value))
	return nil
}

func (c *Compiler) compileArithmetic(node *ast.BinaryExpression) error {
	returnType := c.typeMap[node]
	leftType := c.typeMap[node.Left]
	rightType := c.typeMap[node.Right]

	err := c.Compile(node.Left)
	if err != nil {
		return err
	}
	if leftType == types.IntType && returnType == types.FloatType {
		c.e.ToFloat()
	}
	err = c.Compile(node.Right)
	if err != nil {
		return err
	}
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

	return nil
}
func resolveType(left, right *types.Type) *types.Type {

	if left == right {
		return left
	}

	return types.FloatType
}
func (c *Compiler) compileComparison(node *ast.BinaryExpression) error {

	operator := node.Operator.Type

	leftType := c.typeMap[node.Left]
	rightType := c.typeMap[node.Right]

	expressionType := resolveType(leftType, rightType)

	err := c.Compile(node.Left)
	if err != nil {
		return err
	}
	if leftType == types.IntType && rightType == types.FloatType {
		c.e.ToFloat()
	}

	err = c.Compile(node.Right)
	if err != nil {
		return err
	}
	if leftType == types.FloatType && rightType == types.IntType {
		c.e.ToFloat()
	}

	switch operator {
	case token.LT:
		if expressionType == types.IntType {
			c.e.LtInt()
		} else {
			c.e.LtFloat()
		}
	case token.LTE:
		if expressionType == types.IntType {
			c.e.LteInt()
		} else {
			c.e.LteFloat()
		}
	case token.GT:
		if expressionType == types.IntType {
			c.e.GtInt()
		} else {
			c.e.GtFloat()
		}
	case token.GTE:
		if expressionType == types.IntType {
			c.e.GteInt()
		} else {
			c.e.GteFloat()
		}
	}

	return nil
}

func (c *Compiler) compileLogical(node *ast.BinaryExpression) error {
	operator := node.Operator.Type

	err := c.Compile(node.Left)
	if err != nil {
		return err
	}
	err = c.Compile(node.Right)
	if err != nil {
		return err
	}

	switch operator {
	case token.AND:
		c.e.AndBool()
	case token.OR:
		c.e.OrBool()
	}

	return nil
}
func (c *Compiler) compileEq(node *ast.BinaryExpression) error {
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}
	err = c.Compile(node.Right)
	if err != nil {
		return err
	}
	c.e.Eq()

	return nil
}

func (c *Compiler) compileNeq(node *ast.BinaryExpression) error {
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}
	err = c.Compile(node.Right)
	if err != nil {
		return err
	}
	c.e.Neq()

	return nil
}
func (c *Compiler) compileConcat(node *ast.BinaryExpression) error {
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}
	err = c.Compile(node.Right)
	if err != nil {
		return err
	}
	c.e.AddString()

	return nil
}

func (c *Compiler) compileBinary(node *ast.BinaryExpression) error {
	operator := node.Operator.Type

	switch operator {
	case token.PLUS, token.MINUS, token.MULTIPLY, token.SLASH:
		return c.compileArithmetic(node)
	case token.GT, token.GTE, token.LT, token.LTE:
		return c.compileComparison(node)
	case token.AND, token.OR:
		return c.compileLogical(node)
	case token.EQ:
		return c.compileEq(node)
	case token.NEQ:
		return c.compileNeq(node)
	case token.CONCAT:
		return c.compileConcat(node)
	default:
		return fmt.Errorf("Can't compile binary operator: %s", operator)
	}

	return nil
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
func (c *Compiler) compileInteger(node *ast.IntegerExpression) error {
	value, err := strconv.Atoi(node.Value)
	if err != nil {
		return fmt.Errorf("Error converting '%s' to integer", node.Value)
	}
	c.e.PushInt(value)
	return nil
}
func (c *Compiler) compileAST(node *ast.AST) error {
	for _, statement := range node.Statements {
		err := c.Compile(statement)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *Compiler) Bytecode() emitter.ByteCode {
	return c.e.Bytecode()
}
