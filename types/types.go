package types

import "strings"

type TypeKind string

type Type struct {
	Kind       TypeKind
	Args       []*Type
	ReturnType *Type
}

var (
	IntType    = &Type{Kind: INT}
	StringType = &Type{Kind: STRING}
	BoolType   = &Type{Kind: BOOL}
	FloatType  = &Type{Kind: FLOAT}
	VoidType   = &Type{Kind: VOID}
	AutoType   = &Type{Kind: AUTO}

	UnknownType = &Type{Kind: UNKNOWN}
)

const (
	INT      TypeKind = "int"
	FLOAT    TypeKind = "float"
	STRING   TypeKind = "string"
	BOOL     TypeKind = "bool"
	FUNCTION TypeKind = "function"

	VOID TypeKind = "void"
	AUTO TypeKind = "auto"

	UNKNOWN TypeKind = "unknown"
	RETURN  TypeKind = "return"
)

func NewFunctionType(args []*Type, ReturnType *Type) *Type {
	t := &Type{
		Kind:       FUNCTION,
		Args:       args,
		ReturnType: ReturnType,
	}

	return t
}

func (t *Type) String() string {
	if t.Kind == INT || t.Kind == STRING || t.Kind == FLOAT || t.Kind == VOID || t.Kind == UNKNOWN || t.Kind == BOOL || t.Kind == AUTO {
		return string(t.Kind)
	}

	if t.Kind == RETURN {
		return string(t.ReturnType.Kind)
	}

	args := []string{}
	for i := range t.Args {
		args = append(args, t.Args[i].String())
	}

	argString := "fn" + "(" + strings.Join(args, ",") + ")"

	elements := []string{argString, t.ReturnType.String()}

	return strings.Join(elements, " ")

}
