package parser

import (
	"fmt"

	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/types"
)

type TypeError error

type TypeChecker struct {
	errors    []TypeError
	variables map[string]*types.Type
}

func NewTypeChcker() *TypeChecker {
	t := &TypeChecker{
		variables: map[string]*types.Type{},
		errors:    []TypeError{},
	}

	return t
}

func (t *TypeChecker) TypeCheck(node ast.Node) *types.Type {
	switch node := node.(type) {
	case *ast.AST:
		return t.typeAST(node)
	case *ast.BlockStatement:
		return t.typeBlockStatement(node)
	case ast.ExpressionStatement:
		return t.TypeCheck(node.Inside)
	case ast.NumberExpression:
		return types.IntType
	case ast.StringExpression:
		return types.StringType
	case ast.BooleanExpression:
		return types.BoolType
	case ast.LetStatement:
		return t.typeLetStatement(node)
	case ast.ReturnStatement:
		return t.typeReturnStatement(node)
	case ast.IfStatement:
		return t.typeIfStatement(node)
	case ast.IdentifierExpression:
		return t.typeIdentifierExpression(node)
	case ast.FunctionStatement:
		return t.typeFunctionStatement(node)
	default:
		t.registerError("Can't check type of '%T'", node)
		return types.UnknownType
	}
}

func (t *TypeChecker) typeFunctionStatement(node ast.FunctionStatement) *types.Type {
	functiontype := &types.Type{Kind: types.FUNCTION}

	if node.ReturnType != nil {
		functiontype.ReturnType = node.ReturnType
	} else {
		functiontype.ReturnType = types.VoidType
	}
	// TODO: Check for return statement and see if it matches the returntype mentioned in function header.

	functiontype.Args = []*types.Type{}

	for _, arg := range node.Type {
		argtype := arg
		functiontype.Args = append(functiontype.Args, argtype)
	}

	return functiontype
}

func (t *TypeChecker) typeIdentifierExpression(node ast.IdentifierExpression) *types.Type {
	atype := t.gettype(node.Value.Value)

	if atype == types.UnknownType {
		t.registerError("Identifier '%s' not declared in this scope.", node.Value.Value)
	}

	return atype
}
func (t *TypeChecker) typeIfStatement(node ast.IfStatement) *types.Type {
	condtype := t.TypeCheck(node.Condition)

	if condtype == types.UnknownType {
		return types.UnknownType
	}

	if condtype != types.BoolType {
		t.registerError("Expected condition to be of type BOOLEAN, got '%s'", condtype.Kind)
		return types.UnknownType
	}

	constype := t.TypeCheck(node.Consequence)
	if constype == types.UnknownType {
		return types.UnknownType
	}

	if node.Alternative != nil {
		altype := t.TypeCheck(node.Alternative)

		if altype == types.UnknownType {
			return types.UnknownType
		}
	}

	return condtype
}
func (t *TypeChecker) typeReturnStatement(node ast.ReturnStatement) *types.Type {
	valuetype := t.TypeCheck(node.Value)

	return valuetype
}
func (t *TypeChecker) typeLetStatement(node ast.LetStatement) *types.Type {
	valuetype := t.TypeCheck(node.Value)

	if node.Type != nil {
		pretype := getType(node.Type.Value)

		if valuetype != pretype {
			t.registerError("Expected type of %s for value in let statement, got %s.", pretype.Kind, valuetype.Kind)
			return types.UnknownType
		}
	}

	return t.settype(node.Name.Value, valuetype)
}
func (t *TypeChecker) typeAST(node *ast.AST) *types.Type {
	tp := types.VoidType
	for _, statement := range node.Statements {
		tp = t.TypeCheck(statement)
		if tp == types.UnknownType {
			return tp
		}
	}
	return tp
}

func (t *TypeChecker) typeBlockStatement(node *ast.BlockStatement) *types.Type {
	tp := types.VoidType
	for _, statement := range node.Statements {
		tp = t.TypeCheck(statement)
		if tp == types.UnknownType {
			return tp
		}
	}
	return tp
}
func (t *TypeChecker) registerError(format string, args ...any) {
	t.errors = append(t.errors, fmt.Errorf(format, args...))
}
func (t *TypeChecker) Errors() []TypeError {
	return t.errors
}

func getType(value string) *types.Type {
	switch value {
	case "int":
		return types.IntType
	case "bool":
		return types.BoolType
	case "string":
		return types.StringType
	default:
		return types.UnknownType
	}
}
func (t *TypeChecker) varexists(name string) bool {
	_, ok := t.variables[name]

	return ok
}
func (t *TypeChecker) gettype(name string) *types.Type {
	if t.varexists(name) {
		return t.variables[name]
	}
	return types.UnknownType
}
func (t *TypeChecker) settype(name string, valuetype *types.Type) *types.Type {
	if t.varexists(name) {
		if valuetype != t.gettype(name) {
			t.registerError("Variable '%s' exists with type: %s. Tried to assign to type %s", name, t.gettype(name).Kind, valuetype.Kind)
			return types.UnknownType
		}
	}
	t.variables[name] = valuetype
	return valuetype
}
