package typechecker

import (
	"fmt"

	"github.com/pspiagicw/tremor/types"
)

type TypeScope struct {
	symbols map[string]*types.Type

	Outer *TypeScope
}

func (t *TypeScope) SetupBuiltinFunctions() {
	t.Add("print", types.NewFunctionType([]*types.Type{types.StringType}, types.VoidType))
}

func (t *TypeScope) Add(name string, nodetype *types.Type) error {
	// TODO: Research into how to bifurcate this thing. Like either merge both into one map or implement more features. (Done)
	if val, ok := t.symbolExists(name); ok {
		return fmt.Errorf("Symbol '%s', already declared with type '%s'", name, val)
	}
	t.symbols[name] = nodetype
	return nil
}

func (t *TypeScope) Get(name string) *types.Type {
	val, ok := t.symbolExists(name)

	if !ok {
		return types.UnknownType
	}

	return val
}

func (t *TypeScope) symbolExists(name string) (*types.Type, bool) {
	val, ok := t.symbols[name]

	if ok {
		return val, ok
	}

	if t.Outer != nil {
		return t.Outer.symbolExists(name)
	}

	return nil, false
}

func NewEnclosedScope(outer *TypeScope) *TypeScope {
	s := &TypeScope{
		symbols: map[string]*types.Type{},
	}

	s.Outer = outer

	return s
}
func NewScope() *TypeScope {
	s := &TypeScope{
		symbols: map[string]*types.Type{},
		Outer:   nil,
	}

	return s
}
