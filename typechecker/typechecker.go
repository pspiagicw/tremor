package typechecker

import (
	"fmt"
	"reflect"

	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/token"
	"github.com/pspiagicw/tremor/types"
)

type TypeError error

type TypeMap map[ast.Node]*types.Type

type TypeChecker struct {
	errors  []TypeError
	info    []string
	typeMap TypeMap
}

func (t *TypeChecker) Flush() {
	t.errors = []TypeError{}
}

func NewTypeChecker() *TypeChecker {
	t := &TypeChecker{
		errors:  []TypeError{},
		info:    []string{},
		typeMap: make(map[ast.Node]*types.Type),
	}

	return t
}

func (t *TypeChecker) Map() TypeMap {
	return t.typeMap
}

func (t *TypeChecker) TypeCheck(node ast.Node, scope *TypeScope) *types.Type {
	nodeType := types.UnknownType
	switch node := node.(type) {
	case *ast.AST:
		nodeType = t.typeAST(node, scope)
	case *ast.BlockStatement:
		nodeType = t.typeBlockStatement(node, scope)
	case *ast.ExpressionStatement:
		nodeType = t.TypeCheck(node.Inside, scope)
	case *ast.IntegerExpression:
		nodeType = types.IntType
	case *ast.FloatExpression:
		nodeType = types.FloatType
	case *ast.StringExpression:
		nodeType = types.StringType
	case *ast.BooleanExpression:
		nodeType = types.BoolType
	case *ast.LetStatement:
		nodeType = t.typeLetStatement(node, scope)
	case *ast.ReturnStatement:
		nodeType = t.typeReturnStatement(node, scope)
	case *ast.IfStatement:
		nodeType = t.typeIfStatement(node, scope)
	case *ast.IdentifierExpression:
		nodeType = t.typeIdentifierExpression(node, scope)
	case *ast.FunctionStatement:
		nodeType = t.typeFunctionStatement(node, scope)
	case *ast.LambdaExpression:
		nodeType = t.typeLambdaExpression(node, scope)
	case *ast.FunctionCallExpression:
		nodeType = t.typeFunctionCall(node, scope)
	case *ast.BinaryExpression:
		nodeType = t.typeBinaryExpression(node, scope)
	case *ast.ParenthesisExpression:
		nodeType = t.typeParenthesisExpression(node, scope)
	case *ast.AssignmentStatement:
		nodeType = t.typeAssignmentExpression(node, scope)
	case *ast.PrefixExpression:
		nodeType = t.typePrefixExpression(node, scope)
	case *ast.ArrayExpression:
		nodeType = t.typeArrayExpression(node, scope)
	case *ast.HashExpression:
		nodeType = t.typeHashExpression(node, scope)
	case *ast.IndexExpression:
		nodeType = t.typeIndexExpression(node, scope)
	case *ast.ClassStatement:
		nodeType = t.typeClassStatement(node, scope)
	default:
		t.registerError("Can't check type of '%T'", node)
		return types.UnknownType
	}

	t.typeMap[node] = nodeType

	return nodeType
}
func (t *TypeChecker) typeClassStatement(node *ast.ClassStatement, scope *TypeScope) *types.Type {
	classType := &types.Type{Kind: types.CLASS}

	err := scope.Add(node.Name.Value, classType)
	if err != nil {
		t.addError(err)
		return types.UnknownType
	}

	return classType
}
func (t *TypeChecker) typeArrayIndex(node *ast.IndexExpression, scope *TypeScope) *types.Type {
	arrayType := t.TypeCheck(node.Caller, scope)

	indexType := t.TypeCheck(node.Index, scope)

	if indexType != types.IntType {
		t.registerError("Can't index array with index type %s, requires int", indexType)
		return types.UnknownType
	}

	return arrayType.KeyType
}
func (t *TypeChecker) typeHashAccess(node *ast.IndexExpression, scope *TypeScope) *types.Type {

	hashType := t.TypeCheck(node.Caller, scope)

	indexType := t.TypeCheck(node.Index, scope)

	if indexType != hashType.KeyType {
		t.registerError("Can't index hash with keyType %s, requires %s", indexType, hashType.KeyType)
		return types.UnknownType
	}

	return hashType.ValueType
}
func (t *TypeChecker) typeIndexExpression(node *ast.IndexExpression, scope *TypeScope) *types.Type {
	hostType := t.TypeCheck(node.Caller, scope)

	switch hostType.Kind {
	case types.ARRAY:
		return t.typeArrayIndex(node, scope)
	case types.HASH:
		return t.typeHashAccess(node, scope)
	default:
		t.registerError("Can't index/access type of %s, requires array/hash", hostType)
		return types.UnknownType
	}
}
func isConcreteType(kind types.TypeKind) bool {
	if kind == types.ARRAY || kind == types.HASH || kind == types.CLASS || kind == types.FUNCTION {
		return false
	}
	return true
}
func (t *TypeChecker) typeHashExpression(node *ast.HashExpression, scope *TypeScope) *types.Type {
	hashType := &types.Type{Kind: types.HASH}

	if len(node.Keys) == 0 {
		return types.VoidType
	}

	// DONE: Add a check that only concrete types can be keys (not arrays or hashes or custom types).
	expectedKeyType := t.TypeCheck(node.Keys[0], scope)
	if expectedKeyType == types.UnknownType {
		return types.UnknownType
	}

	if expectedKeyType == types.VoidType {
		t.registerError("Expected concrete type, got void")
		return types.UnknownType
	}

	if !isConcreteType(expectedKeyType.Kind) {
		t.registerError("Expected concrete type, got %s", expectedKeyType.Kind)
		return types.UnknownType
	}

	expectedValueType := t.TypeCheck(node.Values[0], scope)
	if expectedValueType == types.UnknownType {
		return types.UnknownType
	}

	if expectedValueType == types.VoidType {
		t.registerError("Expected concrete type, got void")
		return types.UnknownType
	}

	for i, key := range node.Keys {
		keyType := t.TypeCheck(key, scope)

		// TODO: Combine unknown and voidtype check
		if keyType == types.UnknownType {
			return types.UnknownType
		}

		if keyType == types.VoidType {
			t.registerError("Expected some concrete type, got void")
		}

		if keyType != expectedKeyType {
			t.registerError("Key of hash not of same type, got %s, expected %s", keyType, expectedKeyType)
		}

		valueType := t.TypeCheck(node.Values[i], scope)

		// TODO: Combine unknown and voidtype check
		if valueType == types.UnknownType {
			return types.UnknownType
		}

		if valueType == types.VoidType {
			t.registerError("Expected some concrete type, got void")
		}

		if valueType != expectedValueType {
			t.registerError("Key of hash not of same type, got %s, expected %s", valueType, expectedValueType)
		}
	}

	hashType.KeyType = expectedKeyType
	hashType.ValueType = expectedValueType

	return hashType
}
func (t *TypeChecker) typeArrayExpression(node *ast.ArrayExpression, scope *TypeScope) *types.Type {
	// TODO: Check this once.
	arrType := &types.Type{Kind: types.ARRAY}

	if len(node.Elements) == 0 {
		return types.VoidType
	}

	expectedType := t.TypeCheck(node.Elements[0], scope)
	if expectedType == types.UnknownType {
		return types.UnknownType
	}

	if expectedType == types.VoidType {
		t.registerError("Expected some concrete type, got void")
		return types.UnknownType
	}

	for _, element := range node.Elements[1:] {
		elementType := t.TypeCheck(element, scope)
		if elementType == types.UnknownType {
			return types.UnknownType
		}

		// TODO: Combine unknown and voidtype check
		if elementType == types.VoidType {
			t.registerError("Expected some concrete type, got void")
			return types.UnknownType
		}

		if elementType != expectedType {
			t.registerError("Elements of array not of same type, got %s, expected %s", elementType, expectedType)
			return types.UnknownType
		}
	}

	arrType.KeyType = expectedType

	return arrType
}
func (t *TypeChecker) typePrefixExpression(node *ast.PrefixExpression, scope *TypeScope) *types.Type {
	nodeType := t.TypeCheck(node.Right, scope)

	if node.Operator.Type == token.MINUS {
		if nodeType != types.IntType && nodeType != types.FloatType {
			t.registerError("Expected type to be int or float, got %s", nodeType)
			return types.UnknownType
		}
		return nodeType
	}

	if node.Operator.Type == token.NOT {
		if nodeType != types.BoolType {
			t.registerError("Expected type to be bool, got %s", nodeType)
			return types.UnknownType
		}
		return nodeType
	}

	return types.VoidType
}
func (t *TypeChecker) typeAssignmentExpression(node *ast.AssignmentStatement, scope *TypeScope) *types.Type {
	valuetype := t.TypeCheck(node.Value, scope)
	// TODO: Check if the value type is void, can't assign void to anything.

	existingType := scope.Get(node.Name.Value)

	if existingType == types.UnknownType {
		t.registerError("Variable %s not declared", node.Name.Value)
		return types.UnknownType
	}

	if !reflect.DeepEqual(existingType, valuetype) {
		t.registerError("Type of expression doesn't match type of declared variable.")
		return types.UnknownType
	}

	return valuetype
}
func (t *TypeChecker) typeParenthesisExpression(node *ast.ParenthesisExpression, scope *TypeScope) *types.Type {
	return t.TypeCheck(node.Inside, scope)
}

func (t *TypeChecker) typeBinaryExpression(node *ast.BinaryExpression, scope *TypeScope) *types.Type {
	left := t.TypeCheck(node.Left, scope)

	if left == types.UnknownType {
		return types.UnknownType
	}

	right := t.TypeCheck(node.Right, scope)
	if right == types.UnknownType {
		return types.UnknownType
	}

	operator := node.Operator.Type

	resolver, ok := binaryResolvers[operator]

	if !ok {
		t.registerError("Unsupported operator: %q", operator)
		return types.UnknownType
	}

	expType, err := resolver(left, right)
	if err != nil {
		t.registerError("%s", err.Error())
		return types.UnknownType
	}

	return expType
}
func (t *TypeChecker) typeFunctionCall(node *ast.FunctionCallExpression, scope *TypeScope) *types.Type {
	ftype := scope.Get(node.Caller.String())

	if ftype == types.UnknownType {
		t.registerError("Function '%s', not declared in this scope.", node.Caller.String())
		return types.UnknownType
	} else if ftype.Kind != types.FUNCTION {
		t.registerError("%s is not a function!", node.Caller.String())
		return types.UnknownType
	}

	// DONE: Add test for function call, test arity etc.
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

func (t *TypeChecker) typeFunctionStatement(node *ast.FunctionStatement, scope *TypeScope) *types.Type {
	functiontype := &types.Type{Kind: types.FUNCTION}

	functiontype.ReturnType = node.ReturnType

	newScope := NewEnclosedScope(scope)

	functiontype.Args = []*types.Type{}

	for i, argtype := range node.Type {
		name := node.Args[i].Value
		functiontype.Args = append(functiontype.Args, argtype)
		newScope.Add(name, argtype)
	}

	// TODO: Check if recursion in typechecker works.
	newScope.Add(node.Name.Value, functiontype)

	bodyType := t.TypeCheck(node.Body, newScope)

	if bodyType.Kind == types.RETURN {
		if bodyType.AlwaysReturns == false {
			t.registerError("Expected to always return, it doesn't.")
			return types.UnknownType
		}
		bodyType = bodyType.ReturnType
	}

	if !reflect.DeepEqual(bodyType, functiontype.ReturnType) {
		t.registerError("Expected return type of %s, got %s", functiontype.ReturnType, bodyType)
		return bodyType
	}

	err := scope.Add(node.Name.Value, functiontype)
	if err != nil {
		t.addError(err)
		return types.UnknownType
	}

	return functiontype
}

func (t *TypeChecker) typeLambdaExpression(node *ast.LambdaExpression, scope *TypeScope) *types.Type {
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

	if bodyType.Kind == types.RETURN {
		if bodyType.AlwaysReturns == false {
			t.registerError("Expected to always return, it doesn't.")
			return types.UnknownType
		}
		bodyType = bodyType.ReturnType
	}

	if bodyType != functiontype.ReturnType {
		t.registerError("Expected return type of %s, got %s", functiontype.ReturnType, bodyType)
		return bodyType
	}
	return functiontype
}

func (t *TypeChecker) typeIdentifierExpression(node *ast.IdentifierExpression, scope *TypeScope) *types.Type {
	atype := scope.Get(node.Value.Value)

	if atype == types.UnknownType {
		t.registerError("Symbol '%s' not declared in this scope.", node.Value.Value)
	}

	return atype
}
func (t *TypeChecker) typeIfStatement(node *ast.IfStatement, scope *TypeScope) *types.Type {
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

	var alttype *types.Type = nil
	if node.Alternative != nil {
		alttype = t.TypeCheck(node.Alternative, scope)

		if alttype == types.UnknownType {
			return types.UnknownType
		}
	}

	if alttype != nil {
		// Both consequence and alternative are present
		if constype.Kind != alttype.Kind {
			t.registerError("Type of consequence different from alternative, %s and %s", constype.Kind, alttype.Kind)
			return types.UnknownType
		}

		if constype.Kind == types.RETURN {
			newReturnType := &types.Type{Kind: types.RETURN}
			newReturnType.AlwaysReturns = constype.AlwaysReturns && alttype.AlwaysReturns
			newReturnType.ReturnType = constype.ReturnType

			return newReturnType
		}

		return constype
	}

	// Conditional returns are not a thing.
	return types.VoidType
}
func (t *TypeChecker) typeReturnStatement(node *ast.ReturnStatement, scope *TypeScope) *types.Type {
	valuetype := t.TypeCheck(node.Value, scope)

	if valuetype == types.UnknownType {
		return valuetype
	}

	rt := &types.Type{Kind: types.RETURN}
	rt.ReturnType = valuetype
	rt.AlwaysReturns = true

	return rt
}
func (t *TypeChecker) typeLetStatement(node *ast.LetStatement, scope *TypeScope) *types.Type {
	valuetype := t.TypeCheck(node.Value, scope)

	// DONE: Check if the value is void (can't assign void to anything.)
	if valuetype == types.VoidType {
		t.registerError("Can't assign void to value.")
		return types.UnknownType
	}

	if valuetype == types.UnknownType {
		t.registerError("Can't assign uknown to value.")
		return types.UnknownType
	}

	pretype := node.Type

	if pretype == types.AutoType {
		t.registerInfo("Auto-typed into %s", valuetype)
		pretype = valuetype
	} else if !reflect.DeepEqual(valuetype, pretype) {
		t.registerError("Expected type of %s (pre-type), got %s", pretype.Kind, valuetype.Kind)
		return types.UnknownType
	}

	err := scope.Add(node.Name.Value, valuetype)
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

		statementType := t.TypeCheck(statement, scope)

		if statementType == types.UnknownType {
			return statementType
		}

		if statementType.AlwaysReturns && statementType.Kind == types.RETURN {
			return statementType
		}

	}
	return types.VoidType
}
func (t *TypeChecker) registerError(format string, args ...any) {
	t.errors = append(t.errors, fmt.Errorf(format, args...))
}
func (t *TypeChecker) registerInfo(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	t.info = append(t.info, msg)
}
func (t *TypeChecker) addError(err error) {
	t.errors = append(t.errors, err)
}
func (t *TypeChecker) Errors() []TypeError {
	return t.errors
}
