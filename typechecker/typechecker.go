package parser

import (
	"fmt"

	"github.com/pspiagicw/fener/ast"
	"github.com/pspiagicw/fener/types"
)

type TypeError error

type TypeChecker struct {
	errors    []TypeError
	variables map[string]types.Type
}

func NewTypeChcker() *TypeChecker {
	t := &TypeChecker{
		variables: map[string]types.Type{},
		errors:    []TypeError{},
	}

	return t
}

func (t *TypeChecker) TypeCheck(node ast.Node) types.Type {
	switch node := node.(type) {
	case *ast.AST:
		return t.typeAST(node)
	case *ast.BlockStatement:
		return t.typeBlockStatement(node)
	case ast.ExpressionStatement:
		return t.TypeCheck(node.Inside)
	case ast.NumberExpression:
		return types.INT
	case ast.StringExpression:
		return types.STRING
	case ast.BooleanExpression:
		return types.BOOL
	case ast.LetStatement:
		return t.typeLetStatement(node)
	case ast.ReturnStatement:
		return t.typeReturnStatement(node)
	case ast.IfStatement:
		return t.typeIfStatement(node)
	case ast.IdentifierExpression:
		return t.typeIdentifierExpression(node)
	default:
		t.registerError("Can't check type of '%T'", node)
		return types.UNKNOWN
	}
}
func (t *TypeChecker) typeIdentifierExpression(node ast.IdentifierExpression) types.Type {
	atype := t.gettype(node.Value.Value)

	if atype == types.UNKNOWN {
		t.registerError("Identifier '%s' not declared in this scope.", node.Value.Value)
	}

	return atype
}
func (t *TypeChecker) typeIfStatement(node ast.IfStatement) types.Type {
	condtype := t.TypeCheck(node.Condition)

	if condtype == types.UNKNOWN {
		return types.UNKNOWN
	}

	if condtype != types.BOOL {
		t.registerError("Expected condition to be of type BOOLEAN, got '%s'", condtype)
		return types.UNKNOWN
	}

	constype := t.TypeCheck(node.Consequence)
	if constype == types.UNKNOWN {
		return types.UNKNOWN
	}

	if node.Alternative != nil {
		altype := t.TypeCheck(node.Alternative)

		if altype == types.UNKNOWN {
			return types.UNKNOWN
		}
	}

	return condtype
}
func (t *TypeChecker) typeReturnStatement(node ast.ReturnStatement) types.Type {
	valuetype := t.TypeCheck(node.Value)

	return valuetype
}
func (t *TypeChecker) typeLetStatement(node ast.LetStatement) types.Type {
	valuetype := t.TypeCheck(node.Value)

	if node.Type != nil {
		pretype := getType(node.Type.Value)

		if valuetype != pretype {
			t.registerError("Expected type of %s for value in let statement, got %s.", pretype, valuetype)
			return types.UNKNOWN
		}
	}

	return t.settype(node.Name.Value, valuetype)
}
func (t *TypeChecker) typeAST(node *ast.AST) types.Type {
	tp := types.VOID
	for _, statement := range node.Statements {
		tp = t.TypeCheck(statement)
		if tp == types.UNKNOWN {
			return tp
		}
	}
	return tp
}

func (t *TypeChecker) typeBlockStatement(node *ast.BlockStatement) types.Type {
	tp := types.VOID
	for _, statement := range node.Statements {
		tp = t.TypeCheck(statement)
		if tp == types.UNKNOWN {
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

func getType(value string) types.Type {
	switch value {
	case "int":
		return types.INT
	case "bool":
		return types.BOOL
	default:
		return types.UNKNOWN
	}
}
func (t *TypeChecker) varexists(name string) bool {
	_, ok := t.variables[name]

	return ok
}
func (t *TypeChecker) gettype(name string) types.Type {
	if t.varexists(name) {
		return t.variables[name]
	}
	return types.UNKNOWN
}
func (t *TypeChecker) settype(name string, valuetype types.Type) types.Type {
	if t.varexists(name) {
		if valuetype != t.gettype(name) {
			t.registerError("Variable '%s' exists with type: %s. Tried to assign to type %s", name, t.gettype(name), valuetype)
			return types.UNKNOWN
		}
	}
	t.variables[name] = valuetype
	return valuetype
}
