package ast

import (
	"strings"

	"github.com/pspiagicw/tremor/token"
	"github.com/pspiagicw/tremor/types"
)

type StringType int

const (
	_ = iota
	SINGLE_QUOTED
	DOUBLE_QUOTED
	MULTILINE
)

type Node interface {
	String() string
	TypeInfo() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type AST struct {
	Statements []Statement
}

func (a *AST) statementNode() {
}
func (a *AST) TypeInfo() string {
	return "ast"
}

func (a *AST) String() string {

	statementStrings := []string{}

	for _, statement := range a.Statements {
		statementStrings = append(statementStrings, statement.String())
	}

	return strings.Join(statementStrings, " ")
}

type LetStatement struct {
	Name  *token.Token
	Value Expression
	Type  *types.Type
}

func (l *LetStatement) statementNode() {}
func (l *LetStatement) TypeInfo() string {
	return "let-statement"
}
func (l *LetStatement) String() string {
	elements := []string{}

	if l.Type == nil {
		elements = []string{"let", l.Name.Value, "=", l.Value.String()}
	} else {
		elements = []string{"let", l.Name.Value, l.Type.String(), "=", l.Value.String()}
	}

	return strings.Join(elements, " ")
}

type AssignmentStatement struct {
	Name  *token.Token
	Value Expression
}

func (a *AssignmentStatement) TypeInfo() string {
	return "assignment-statement"
}
func (a *AssignmentStatement) expressionNode() {}
func (a *AssignmentStatement) String() string {
	elements := []string{a.Name.Value, "=", a.Value.String()}

	return strings.Join(elements, " ")
}

type IntegerExpression struct {
	Value string
}

func (n *IntegerExpression) TypeInfo() string {
	return "integer-expression"
}
func (n *IntegerExpression) expressionNode() {}
func (n *IntegerExpression) String() string {
	return n.Value
}

type FloatExpression struct {
	Value string
}

func (f *FloatExpression) TypeInfo() string {
	return "float-expression"
}
func (f *FloatExpression) expressionNode() {}
func (f *FloatExpression) String() string {
	return f.Value
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator *token.Token
	Type     *types.Type
}

func (b *BinaryExpression) TypeInfo() string {
	return "binary-expression"
}
func (b *BinaryExpression) expressionNode() {}
func (b *BinaryExpression) String() string {

	elements := []string{b.Left.String(), b.Operator.Value, b.Right.String()}

	return "(" + strings.Join(elements, " ") + ")"
}

type StringExpression struct {
	Value string
	Type  StringType
}

func (s *StringExpression) TypeInfo() string {
	return "string-expression"
}
func (s *StringExpression) expressionNode() {}
func (s *StringExpression) String() string {
	quote := "\""
	endquote := "\""

	switch s.Type {
	case SINGLE_QUOTED:
		quote = "'"
		endquote = "'"
	case MULTILINE:
		quote = "[["
		endquote = "]]"
	}

	return quote + s.Value + endquote
}

type PrefixExpression struct {
	Right    Expression
	Operator *token.Token
}

func (p *PrefixExpression) TypeInfo() string {
	return "prefix-expression"
}
func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) String() string {
	elements := []string{p.Operator.Value, p.Right.String()}

	return "(" + strings.Join(elements, " ") + ")"
}

type BooleanExpression struct {
	Value *token.Token
}

func (b *BooleanExpression) TypeInfo() string {
	return "boolean-expression"
}
func (b *BooleanExpression) expressionNode() {}
func (b *BooleanExpression) String() string {
	return b.Value.Value
}

type ParenthesisExpression struct {
	Inside Expression
}

func (p *ParenthesisExpression) TypeInfo() string {
	return "parenthesis-expression"
}
func (p *ParenthesisExpression) expressionNode() {}
func (p *ParenthesisExpression) String() string {
	return p.Inside.String()
}

type IdentifierExpression struct {
	Value *token.Token
}

func (i *IdentifierExpression) TypeInfo() string {
	return "index-expression"
}
func (i *IdentifierExpression) expressionNode() {}
func (i *IdentifierExpression) String() string {
	return i.Value.Value
}

type FunctionCallExpression struct {
	Caller    Expression
	Arguments []Expression
}

func (f *FunctionCallExpression) TypeInfo() string {
	return "function-call-expression"
}
func (f *FunctionCallExpression) expressionNode() {}
func (f *FunctionCallExpression) String() string {

	args := []string{}

	for _, argument := range f.Arguments {
		args = append(args, argument.String())
	}

	elements := []string{
		f.Caller.String(),
		"(",
		strings.Join(args, ", "),
		")",
	}

	return strings.Join(elements, "")
}

type IndexExpression struct {
	Caller Expression
	Index  Expression
}

func (i *IndexExpression) TypeInfo() string {
	return "index-expression"
}
func (i *IndexExpression) expressionNode() {}
func (i *IndexExpression) String() string {
	elements := []string{i.Caller.String(), "[", i.Index.String(), "]"}

	return strings.Join(elements, "")
}

type FieldExpression struct {
	Caller Expression
	Field  Expression
}

func (f *FieldExpression) TypeInfo() string {
	return "field-expression"
}
func (f *FieldExpression) expressionNode() {}
func (f *FieldExpression) String() string {
	elements := []string{f.Caller.String(), ".", f.Field.String()}

	return strings.Join(elements, "")
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) TypeInfo() string {
	return "block-statement"
}
func (b *BlockStatement) statementNode() {}
func (b *BlockStatement) String() string {
	statementStrings := []string{}

	for _, statement := range b.Statements {
		statementStrings = append(statementStrings, statement.String())
	}

	return strings.Join(statementStrings, " ")
}

type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfStatement) TypeInfo() string {
	return "if-statement"
}
func (i *IfStatement) statementNode() {}
func (i *IfStatement) String() string {
	elements := []string{"if", i.Condition.String(), "then", i.Consequence.String()}

	if i.Alternative != nil {
		elements = append(elements, "else")
		elements = append(elements, i.Alternative.String())
	}

	elements = append(elements, "end")

	return strings.Join(elements, " ")
}

type ExpressionStatement struct {
	Inside Expression
}

func (e *ExpressionStatement) TypeInfo() string {
	return "expression-statement"
}
func (e *ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) String() string {
	return e.Inside.String()
}

type ReturnStatement struct {
	Value Expression
}

func (r *ReturnStatement) TypeInfo() string {
	return "return-statement"
}
func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) String() string {
	elements := []string{"return", r.Value.String()}

	return strings.Join(elements, " ")
}

type FunctionStatement struct {
	Name       *token.Token
	Args       []*token.Token
	Type       []*types.Type
	Body       *BlockStatement
	ReturnType *types.Type
}

func (f *FunctionStatement) TypeInfo() string {
	return "function-statement"
}
func (f *FunctionStatement) statementNode() {}
func (f *FunctionStatement) String() string {
	args := []string{}

	for i, arg := range f.Args {
		name := arg.Value + " " + f.Type[i].String()
		args = append(args, name)
	}

	headerString := f.Name.Value + "(" + strings.Join(args, ", ") + ")"

	elements := []string{}

	if f.ReturnType == nil {
		elements = []string{"fn", headerString, "then", f.Body.String(), "end"}
	} else {
		elements = []string{"fn", headerString, f.ReturnType.String(), "then", f.Body.String(), "end"}
	}

	return strings.Join(elements, " ")
}

type LambdaExpression struct {
	Args       []*token.Token
	Type       []*types.Type
	Body       *BlockStatement
	ReturnType *types.Type
}

func (l *LambdaExpression) TypeInfo() string {
	return "lambad-expression"
}
func (l *LambdaExpression) expressionNode() {}
func (l *LambdaExpression) String() string {
	args := []string{}

	for i, arg := range l.Args {
		name := arg.Value + " " + l.Type[i].String()
		args = append(args, name)
	}

	headerString := "fn(" + strings.Join(args, ", ") + ")"

	elements := []string{}

	if l.ReturnType == nil {
		elements = []string{headerString, "then", l.Body.String(), "end"}
	} else {
		elements = []string{headerString, l.ReturnType.String(), "then", l.Body.String(), "end"}
	}

	return strings.Join(elements, " ")
}

type ArrayExpression struct {
	Elements []Expression
}

func (a *ArrayExpression) TypeInfo() string {
	return "array-expression"
}

func (a *ArrayExpression) expressionNode() {}
func (a *ArrayExpression) String() string {
	args := []string{}

	for _, arg := range a.Elements {
		args = append(args, arg.String())
	}

	elements := []string{"[", strings.Join(args, ", "), "]"}

	return strings.Join(elements, "")
}

type HashExpression struct {
	Keys   []Expression
	Values []Expression
}

func (h *HashExpression) TypeInfo() string {
	return "hash-expression"
}
func (h *HashExpression) expressionNode() {}
func (h *HashExpression) String() string {
	// TODO: Make hash expression have string

	pairs := []string{}

	for i, key := range h.Keys {
		value := h.Values[i]

		combined := key.String() + ": " + value.String()

		pairs = append(pairs, combined)
	}

	elements := []string{"{", strings.Join(pairs, ", "), "}"}

	return strings.Join(elements, "")
}

type ClassStatement struct {
	Name    *token.Token
	Methods []*FunctionStatement
}

func (c *ClassStatement) TypeInfo() string {
	return "class-statement"
}
func (c *ClassStatement) statementNode() {}
func (c *ClassStatement) String() string {

	methodStrings := []string{}

	for _, method := range c.Methods {
		methodStrings = append(methodStrings, method.String())
	}

	methodCombined := strings.Join(methodStrings, " ")

	var elements []string
	if len(methodStrings) == 0 {
		elements = []string{"class", c.Name.Value, "end"}
	} else {
		elements = []string{"class", c.Name.Value, methodCombined, "end"}
	}

	return strings.Join(elements, " ")
}
