package typechecker

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

func NewTypeChecker() *TypeChecker {
	t := &TypeChecker{
		variables: map[string]*types.Type{},
		errors:    []TypeError{},
	}

	// t.addFunc("print", types.NewFunctionType([]*types.Type{types.StringType}, types.VoidType))

	return t
}

//	func (t *TypeChecker) addFunc(name string, returntype *types.Type) {
//		t.settype(name, returntype)
//	}
func (t *TypeChecker) TypeCheck(node ast.Node, scope *TypeScope) *types.Type {
	switch node := node.(type) {
	case *ast.AST:
		return t.typeAST(node, scope)
	case *ast.BlockStatement:
		return t.typeBlockStatement(node, scope)
	case ast.ExpressionStatement:
		return t.TypeCheck(node.Inside, scope)
	case ast.NumberExpression:
		return types.IntType
	case ast.StringExpression:
		return types.StringType
	case ast.BooleanExpression:
		return types.BoolType
	case ast.LetStatement:
		return t.typeLetStatement(node, scope)
	case ast.ReturnStatement:
		return t.typeReturnStatement(node, scope)
	case ast.IfStatement:
		return t.typeIfStatement(node, scope)
	case ast.IdentifierExpression:
		return t.typeIdentifierExpression(node, scope)
	case ast.FunctionStatement:
		return t.typeFunctionStatement(node, scope)
	case ast.FunctionCallExpression:
		return t.typeFunctionCall(node, scope)
	case ast.BinaryExpression:
		return t.typeBinaryExpression(node, scope)
	default:
		t.registerError("Can't check type of '%T'", node)
		return types.UnknownType
	}
}
func (t *TypeChecker) typeBinaryExpression(node ast.BinaryExpression, scope *TypeScope) *types.Type {
	lefttype := t.TypeCheck(node.Left, scope)

	if lefttype == types.UnknownType {
		return types.UnknownType
	}

	righttype := t.TypeCheck(node.Right, scope)
	if righttype == types.UnknownType {
		return types.UnknownType
	}

	if lefttype != righttype {
		t.registerError("Expected same type for both operands, got %s and %s", lefttype, righttype)
	}

	return lefttype
}
func (t *TypeChecker) typeFunctionCall(node ast.FunctionCallExpression, scope *TypeScope) *types.Type {
	ftype := scope.GetFunction(node.Caller.String())

	if ftype == types.UnknownType {
		t.registerError("Function '%s', not declared in this scope.", node.Caller.String())
		return types.UnknownType
	} else if ftype.Kind != types.FUNCTION {
		t.registerError("%s is not a function!", node.Caller.String())
		return types.UnknownType
	}

	if len(ftype.Args) != len(node.Arguments) {
		t.registerError("Function needs %d arguments, got %d", len(ftype.Args), len(node.Arguments))
		return types.UnknownType
	}

	for i, argtype := range ftype.Args {
		actualtype := t.TypeCheck(node.Arguments[i], scope)
		if actualtype != argtype {
			t.registerError("[%d] Function needs argument of type %s, got %s", i, argtype, actualtype)
			return types.UnknownType
		}
	}

	return ftype.ReturnType
}

func (t *TypeChecker) typeFunctionStatement(node ast.FunctionStatement, scope *TypeScope) *types.Type {
	functiontype := &types.Type{Kind: types.FUNCTION}

	functiontype.ReturnType = node.ReturnType

	newScope := NewEnclosedScope(scope)

	functiontype.Args = []*types.Type{}

	for i, argtype := range node.Type {
		name := node.Args[i].Value
		functiontype.Args = append(functiontype.Args, argtype)
		newScope.Add(name, argtype)
	}

	bodyType := t.TypeCheck(node.Body, newScope)

	if bodyType != functiontype.ReturnType {
		t.registerError("Expected return type of %s, got %s", functiontype.ReturnType, bodyType)
		return bodyType
	}
	// TODO: Check for return statement and see if it matches the returntype mentioned in function header. (completed)

	err := scope.AddFunc(node.Name.Value, functiontype)
	if err != nil {
		t.addError(err)
		return types.UnknownType
	}

	return functiontype
}

func (t *TypeChecker) typeIdentifierExpression(node ast.IdentifierExpression, scope *TypeScope) *types.Type {
	atype := scope.GetVariables(node.Value.Value)

	if atype == types.UnknownType {
		t.registerError("Identifier '%s' not declared in this scope.", node.Value.Value)
	}

	return atype
}
func (t *TypeChecker) typeIfStatement(node ast.IfStatement, scope *TypeScope) *types.Type {
	condtype := t.TypeCheck(node.Condition, scope)

	if condtype == types.UnknownType {
		return types.UnknownType
	}

	if condtype != types.BoolType {
		t.registerError("Expected condition to be of type BOOLEAN, got '%s'", condtype.Kind)
		return types.UnknownType
	}

	constype := t.TypeCheck(node.Consequence, scope)
	if constype == types.UnknownType {
		return types.UnknownType
	}

	if node.Alternative != nil {
		altype := t.TypeCheck(node.Alternative, scope)

		if altype == types.UnknownType {
			return types.UnknownType
		}
	}

	return condtype
}
func (t *TypeChecker) typeReturnStatement(node ast.ReturnStatement, scope *TypeScope) *types.Type {
	valuetype := t.TypeCheck(node.Value, scope)

	rt := &types.Type{Kind: types.RETURN}
	rt.ReturnType = valuetype

	return rt
}
func (t *TypeChecker) typeLetStatement(node ast.LetStatement, scope *TypeScope) *types.Type {
	valuetype := t.TypeCheck(node.Value, scope)

	pretype := node.Type

	if valuetype != pretype {
		t.registerError("Expected type of %s for value in let statement, got %s.", pretype.Kind, valuetype.Kind)
		return types.UnknownType
	}

	err := scope.AddVariable(node.Name.Value, valuetype)
	if err != nil {
		t.addError(err)
		return types.UnknownType
	}

	return valuetype

}
func (t *TypeChecker) typeAST(node *ast.AST, scope *TypeScope) *types.Type {
	tp := types.VoidType
	for _, statement := range node.Statements {
		tp = t.TypeCheck(statement, scope)
		if tp == types.UnknownType {
			return tp
		}
	}
	return tp
}

func (t *TypeChecker) typeBlockStatement(node *ast.BlockStatement, scope *TypeScope) *types.Type {
	for _, statement := range node.Statements {
		bodytype := t.TypeCheck(statement, scope)
		if bodytype.Kind == types.RETURN {
			return bodytype.ReturnType
		} else if bodytype == types.UnknownType {
			return bodytype
		}
	}
	return types.VoidType
}
func (t *TypeChecker) registerError(format string, args ...any) {
	t.errors = append(t.errors, fmt.Errorf(format, args...))
}
func (t *TypeChecker) addError(err error) {
	t.errors = append(t.errors, err)
}
func (t *TypeChecker) Errors() []TypeError {
	return t.errors
}

// func getType(value string) *types.Type {
// 	switch value {
// 	case "int":
// 		return types.IntType
// 	case "bool":
// 		return types.BoolType
// 	case "string":
// 		return types.StringType
// 	case "void":
// 		return types.VoidType
// 	default:
// 		return types.UnknownType
// 	}
// }

// func (t *TypeChecker) varexists(name string) bool {
// 	_, ok := t.variables[name]
//
// 	return ok
// }
// func (t *TypeChecker) gettype(name string) *types.Type {
// 	if t.varexists(name) {
// 		return t.variables[name]
// 	}
// 	return types.UnknownType
// }
// func (t *TypeChecker) settype(name string, valuetype *types.Type) *types.Type {
// 	if t.varexists(name) {
// 		if valuetype != t.gettype(name) {
// 			t.registerError("Variable '%s' exists with type: %s. Tried to assign to type %s", name, t.gettype(name).Kind, valuetype.Kind)
// 			return types.UnknownType
// 		}
// 	}
// 	t.variables[name] = valuetype
// 	return valuetype
// }
