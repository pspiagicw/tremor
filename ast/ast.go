package ast

import (
	"strings"

	"github.com/pspiagicw/fener/token"
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

func (a AST) statementNode() {
}

func (a AST) String() string {

	statementStrings := []string{}

	for _, statement := range a.Statements {
		statementStrings = append(statementStrings, statement.String())
	}

	return strings.Join(statementStrings, " ")
}

type LetStatement struct {
	Name  string
	Value Expression
}

func (l LetStatement) statementNode() {}
func (l LetStatement) String() string {
	elements := []string{"let", l.Name, "=", l.Value.String()}

	return strings.Join(elements, " ")
}

type NumberExpression struct {
	Value string
}

func (n NumberExpression) expressionNode() {}
func (n NumberExpression) String() string {
	return n.Value
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator *token.Token
}

func (b BinaryExpression) expressionNode() {}
func (b BinaryExpression) String() string {

	elements := []string{b.Left.String(), b.Operator.Value, b.Right.String()}

	return "(" + strings.Join(elements, " ") + ")"
}

type StringExpression struct {
	Value string
	Type  StringType
}

func (s StringExpression) expressionNode() {}
func (s StringExpression) String() string {
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

func (p PrefixExpression) expressionNode() {}
func (p PrefixExpression) String() string {
	elements := []string{p.Operator.Value, p.Right.String()}

	return "(" + strings.Join(elements, " ") + ")"
}

type BooleanExpression struct {
	Value *token.Token
}

func (b BooleanExpression) expressionNode() {}
func (b BooleanExpression) String() string {
	return b.Value.Value
}

type ParenthesisExpression struct {
	Inside Expression
}

func (p ParenthesisExpression) expressionNode() {}
func (p ParenthesisExpression) String() string {
	return p.Inside.String()
}

type IdentifierExpression struct {
	Value *token.Token
}

func (i IdentifierExpression) expressionNode() {}
func (i IdentifierExpression) String() string {
	return i.Value.Value
}

type FunctionCallExpression struct {
	Caller    Expression
	Arguments []Expression
}

func (f FunctionCallExpression) expressionNode() {}
func (f FunctionCallExpression) String() string {

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

func (i IndexExpression) expressionNode() {}
func (i IndexExpression) String() string {
	elements := []string{i.Caller.String(), "[", i.Index.String(), "]"}

	return strings.Join(elements, "")
}

type FieldExpression struct {
	Caller Expression
	Field  Expression
}

func (f FieldExpression) expressionNode() {}
func (f FieldExpression) String() string {
	elements := []string{f.Caller.String(), ".", f.Field.String()}

	return strings.Join(elements, "")
}

type BlockStatement struct {
	Statements []Statement
}

func (b BlockStatement) statementNode() {}
func (b BlockStatement) String() string {
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

func (i IfStatement) statementNode() {}
func (i IfStatement) String() string {
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

func (e ExpressionStatement) statementNode() {}
func (e ExpressionStatement) String() string {
	return e.Inside.String()
}

type ReturnStatement struct {
	Value Expression
}

func (r ReturnStatement) statementNode() {}
func (r ReturnStatement) String() string {
	elements := []string{"return", r.Value.String()}

	return strings.Join(elements, " ")
}

type FunctionStatement struct {
	Name *token.Token
	Args []*token.Token
	Body *BlockStatement
}

func (f FunctionStatement) statementNode() {}
func (f FunctionStatement) String() string {
	args := []string{}
	for _, arg := range f.Args {
		args = append(args, arg.Value)
	}

	headerString := f.Name.Value + "(" + strings.Join(args, ", ") + ")"

	elements := []string{"fn", headerString, "then", f.Body.String(), "end"}

	return strings.Join(elements, " ")
}
