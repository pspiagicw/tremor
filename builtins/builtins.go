package builtins

import (
	"fmt"
	"os"

	"github.com/pspiagicw/fenc/object"
	"github.com/pspiagicw/tremor/types"
)

type BuiltinDefinition struct {
	Name  string
	BType types.Type
	Impl  func(...object.Object) object.Object
}

var Builtins = []BuiltinDefinition{
	{
		// TODO: Evaluate object system, do we need string() and content() methods, do we need more methods?
		Name:  "print",
		BType: *types.NewFunctionType([]*types.Type{types.NewAnyType([]*types.Type{types.StringType, types.ArrayType, types.HashType, types.BoolType, types.IntType, types.FloatType})}, types.VoidType),
		Impl: func(args ...object.Object) object.Object {
			for _, o := range args {
				fmt.Println(o.String())
			}
			return object.Null{}
		},
	},
	{
		Name:  "len",
		BType: *types.NewFunctionType([]*types.Type{types.NewAnyType([]*types.Type{types.ArrayType, types.HashType, types.StringType})}, types.IntType),
		Impl: func(args ...object.Object) object.Object {
			arg := args[0]

			// DONE: Remove after implementing sub-type checking.
			// DONE: Implement some error object if needed (not needed if implementing sub-type checking)
			switch arg := arg.(type) {
			case object.String:
				return object.CreateInt(len(arg.Value))
			case object.Array:
				return object.CreateInt(len(arg.Values))
			case object.Hash:
				return object.CreateInt(len(arg.Values))
			default:
				return object.Null{}
			}
		},
	},
	{
		Name:  "str",
		BType: *types.NewFunctionType([]*types.Type{types.AnyType}, types.StringType),
		Impl: func(args ...object.Object) object.Object {
			arg := args[0]
			return object.CreateString(arg.Content())
		},
	},
	{
		Name:  "type",
		BType: *types.NewFunctionType([]*types.Type{types.AnyType}, types.StringType),
		Impl: func(args ...object.Object) object.Object {
			arg := args[0]
			// TODO: Find out how to provide the exact type, we can only provide type kind I think.
			return object.CreateString(string(arg.Type()))
		},
	},
	{
		Name:  "exit",
		BType: *types.NewFunctionType([]*types.Type{}, types.VoidType),
		Impl: func(args ...object.Object) object.Object {
			os.Exit(0)
			return object.Null{}
		},
	},
}

func GetBuiltins() map[string]object.Builtin {
	result := make(map[string]object.Builtin, len(Builtins))

	for _, builtin := range Builtins {
		result[builtin.Name] = object.Builtin{
			Internal: builtin.Impl,
		}
	}

	return result
}

// DONE: Add dynamic types (any-type)
// TODO: Add multi-arg support (var-arg)
// DONE: Check if arity is every checked
// Add other builtins

// DONE: print
// DONE: len
// DONE: type (most important)
// TODO: push
// TODO: pop
// TODO: sqrt
// TODO: exp
// TODO: min
// TODO: min/max
// DONE: string (str)
// TODO: lowercase
// TODO: uppercase
// TODO: count
// TODO: exit
