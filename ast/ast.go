package ast

import (
	"github.com/pspiagicw/fener/token"
	"strings"
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

	return strings.Join(elements, " ")
}
