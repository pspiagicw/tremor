package types

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
	VoidType   = &Type{Kind: VOID}

	UnknownType = &Type{Kind: UNKNOWN}
)

const (
	INT      TypeKind = "int"
	FLOAT    TypeKind = "float"
	STRING   TypeKind = "string"
	BOOL     TypeKind = "bool"
	FUNCTION TypeKind = "function"

	VOID TypeKind = "void"

	UNKNOWN TypeKind = "unknown"
)

func NewFunctionType(args []*Type, ReturnType *Type) *Type {
	t := &Type{
		Kind:       FUNCTION,
		Args:       args,
		ReturnType: ReturnType,
	}

	return t
}
