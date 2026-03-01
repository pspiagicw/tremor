package ast

import "github.com/pspiagicw/tremor/token"

func NodeToken(node Node) *token.Token {
	switch n := node.(type) {
	case *LetStatement:
		return n.Name
	case *AssignmentStatement:
		return n.Name
	case *BinaryExpression:
		return n.Operator
	case *PrefixExpression:
		return n.Operator
	case *BooleanExpression:
		return n.Value
	case *IdentifierExpression:
		return n.Value
	case *FunctionStatement:
		return n.Name
	case *ClassStatement:
		return n.Name
	case *ExpressionStatement:
		if n.Inside != nil {
			return NodeToken(n.Inside)
		}
	case *ReturnStatement:
		if n.Value != nil {
			return NodeToken(n.Value)
		}
	case *IfStatement:
		if n.Condition != nil {
			return NodeToken(n.Condition)
		}
	case *ParenthesisExpression:
		if n.Inside != nil {
			return NodeToken(n.Inside)
		}
	case *FunctionCallExpression:
		if n.Caller != nil {
			return NodeToken(n.Caller)
		}
	case *IndexExpression:
		if n.Caller != nil {
			return NodeToken(n.Caller)
		}
	case *FieldExpression:
		if n.Caller != nil {
			return NodeToken(n.Caller)
		}
	case *ArrayExpression:
		if len(n.Elements) > 0 {
			return NodeToken(n.Elements[0])
		}
	case *HashExpression:
		if len(n.Keys) > 0 {
			return NodeToken(n.Keys[0])
		}
	case *BlockStatement:
		if len(n.Statements) > 0 {
			return NodeToken(n.Statements[0])
		}
	case *AST:
		if len(n.Statements) > 0 {
			return NodeToken(n.Statements[0])
		}
	}
	return nil
}
