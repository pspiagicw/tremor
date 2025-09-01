package typechecker

import (
	"fmt"

	"github.com/pspiagicw/tremor/types"
)

type TypeScope struct {
	variables map[string]*types.Type

	functions map[string]*types.Type

	Outer *TypeScope
}

func (t *TypeScope) SetupBuiltinFunctions() {
	t.AddFunc("print", types.NewFunctionType([]*types.Type{types.StringType}, types.VoidType))
}

func (t *TypeScope) AddVariable(name string, vtype *types.Type) error {
	if val, ok := t.varExists(name); ok {
		return fmt.Errorf("Variable '%s', already declared with type '%s'", name, val)
	}

	t.variables[name] = vtype

	return nil
}

func (t *TypeScope) Add(name string, nodetype *types.Type) {
	// TODO: Research into how to bifurcate this thing. Like either merge both into one map or implement more features.
	if nodetype.Kind == types.FUNCTION {
		t.AddFunc(name, nodetype)
	}
	t.AddVariable(name, nodetype)
}

func (t *TypeScope) GetFunction(name string) *types.Type {
	val, ok := t.funcExists(name)

	if !ok {
		return types.UnknownType
	}

	return val
}
func (t *TypeScope) GetVariables(name string) *types.Type {
	val, ok := t.varExists(name)

	if !ok {
		return types.UnknownType
	}

	return val
}

func (t *TypeScope) varExists(name string) (*types.Type, bool) {
	val, ok := t.variables[name]

	if ok {
		return val, ok
	}

	if t.Outer != nil {
		return t.Outer.varExists(name)
	}

	return nil, false
}

func (t *TypeScope) funcExists(name string) (*types.Type, bool) {
	val, ok := t.functions[name]

	if ok {
		return val, ok
	}

	if t.Outer != nil {
		return t.Outer.funcExists(name)
	}

	return nil, false

}

func (t *TypeScope) AddFunc(name string, functype *types.Type) error {
	if val, ok := t.funcExists(name); ok {
		return fmt.Errorf("Function '%s', already declared with type '%s'", name, val)
	}

	t.functions[name] = functype

	return nil
}

func NewEnclosedScope(outer *TypeScope) *TypeScope {
	s := &TypeScope{
		variables: map[string]*types.Type{},
		functions: map[string]*types.Type{},
	}

	s.Outer = outer

	return s
}
func NewScope() *TypeScope {
	s := &TypeScope{
		variables: map[string]*types.Type{},
		functions: map[string]*types.Type{},
		Outer:     nil,
	}

	return s
}
