package builtins

import (
	"fmt"

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
		Name:  "print",
		BType: *types.NewFunctionType([]*types.Type{types.StringType}, types.VoidType),
		Impl: func(args ...object.Object) object.Object {
			for _, o := range args {
				fmt.Println(o.Content())
			}
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

// Add other builtins

// TODO: print
// TODO: len
// TODO: push
// TODO: pop
// TODO: sqrt
// TODO: exp
// TODO: min
// TODO: min/max
// TODO: string
// TODO: lowercase
// TODO: uppercase
// TODO: count
// TODO: type (most important)
// TODO: exit
