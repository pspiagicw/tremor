package parser

import (
	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/token"
	"github.com/pspiagicw/tremor/types"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParseFnMap[p.current.Type]

	if prefixFn == nil {
		p.registerError(FAILED_PREFIX_MESSAGE, p.current.Type)
		return nil
	}

	if p.peek.Type == token.ASSIGN {
		// Exceptional case, the expression is actually a assignment statement
		return p.parseAssignmentStatement()
	}

	// Should have a advance function automatically in it.
	left := prefixFn()

	for p.current.Type != token.EOF && precedence < p.currentPrecedence() {

		infixFn := p.infixParseFnMap[p.current.Type]

		if infixFn == nil {
			return left
		}

		left = infixFn(left)

		if left == nil {
			return nil
		}

	}

	return left
}

func (p *Parser) currentPrecedence() int {
	val, ok := precedence[p.current.Type]

	if ok {
		return val
	}

	return LOWEST
}

func (p *Parser) parseIntegerExpression() ast.Expression {
	n := &ast.IntegerExpression{
		Value: p.current.Value,
	}

	p.advance()
	return n
}
func (p *Parser) parseFloatExpression() ast.Expression {
	n := &ast.FloatExpression{
		Value: p.current.Value,
	}

	p.advance()
	return n
}
func (p *Parser) parseStringExpression() ast.Expression {
	s := &ast.StringExpression{
		Value: p.current.Value,
	}

	switch p.current.Type {
	case token.STRING_DOUBLE:
		s.Type = ast.DOUBLE_QUOTED
	case token.STRING_SINGLE:
		s.Type = ast.SINGLE_QUOTED
	default:
		s.Type = ast.MULTILINE
	}

	p.advance()
	return s
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	i := &ast.IdentifierExpression{
		Value: p.current,
	}

	p.advance()
	return i
}
func (p *Parser) parseBooleanExpression() ast.Expression {
	b := &ast.BooleanExpression{
		Value: p.current,
	}
	p.advance()

	return b
}
func (p *Parser) parseParenthesisExpression() ast.Expression {
	exp := &ast.ParenthesisExpression{}

	p.advance()

	exp.Inside = p.parseExpression(LOWEST)

	// Skip over the ending round brackets.
	p.advance()

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Operator: p.current,
	}

	p.advance()
	// Hard coded precedence value!
	exp.Right = p.parseExpression(UNARY)

	return exp
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	operator := p.current

	b := &ast.BinaryExpression{
		Left:     left,
		Operator: operator,
	}

	operatorPrecedence := p.currentPrecedence()

	p.advance()

	b.Right = p.parseExpression(operatorPrecedence)

	return b
}
func (p *Parser) parseFunctionCallExpression(left ast.Expression) ast.Expression {
	f := &ast.FunctionCallExpression{
		Caller: left,
	}

	p.advance()

	f.Arguments = []ast.Expression{}

	for p.current.Type != token.RPAREN {
		f.Arguments = append(f.Arguments, p.parseExpression(LOWEST))

		if p.current.Type == token.COMMA {
			p.advance()
		}
	}

	// TODO: There should be a expect right here right ?
	// TODO: Check all parse functions.
	p.advance()

	return f
}
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	i := &ast.IndexExpression{
		Caller: left,
	}

	p.advance()

	i.Index = p.parseExpression(LOWEST)

	p.advance()

	return i
}
func (p *Parser) parseFieldExpression(left ast.Expression) ast.Expression {
	f := &ast.FieldExpression{
		Caller: left,
	}

	p.advance()

	// Hard coded dot precedence
	f.Field = p.parseExpression(FIELD)

	return f
}

func (p *Parser) parseLambdaExpression() ast.Expression {
	p.advance() // remove the fn token

	l := &ast.LambdaExpression{}

	p.expect(token.LPAREN)

	l.Args = []*token.Token{}
	l.Type = []*types.Type{}

	for p.current.Type != token.RPAREN {
		arg := p.expect(token.IDENTIFIER)
		l.Args = append(l.Args, arg)

		argtype := p.parseTypeDec(false)
		l.Type = append(l.Type, argtype)

		if p.current.Type == token.RPAREN {
			break
		} else if p.current.Type == token.COMMA {
			p.advance()
		} else {
			p.registerError(FAILED_FUNCTION_MESSAGE, p.current.Type)
		}
	}

	p.expect(token.RPAREN)

	l.ReturnType = p.parseTypeDec(true)

	if l.ReturnType == types.AutoType {
		l.ReturnType = types.VoidType
	}

	p.expect(token.THEN)

	l.Body = p.parseBlockStatement()

	p.expect(token.END)

	return l

}

func (p *Parser) parseArrayExpression() ast.Expression {
	p.advance() // Move ahead the [

	a := &ast.ArrayExpression{
		Elements: []ast.Expression{},
	}

	for p.current.Type != token.RSQUARE {
		element := p.parseExpression(LOWEST)
		a.Elements = append(a.Elements, element)
		if p.current.Type == token.RSQUARE {
			break
		} else if p.current.Type == token.COMMA {
			p.advance()
		} else {
			// TODO: Add a different message
			p.registerError(FAILED_FUNCTION_MESSAGE, p.current.Type)
		}
	}

	p.expect(token.RSQUARE)

	return a
}
func (p *Parser) parseHashExpression() ast.Expression {
	p.advance() // Move over the {

	a := &ast.HashExpression{
		Keys:   []ast.Expression{},
		Values: []ast.Expression{},
	}

	for p.current.Type != token.RBRACE {
		key := p.parseExpression(LOWEST)
		p.expect(token.COLON)
		value := p.parseExpression(LOWEST)

		a.Keys = append(a.Keys, key)
		a.Values = append(a.Values, value)

		if p.current.Type == token.RBRACE {
			break
		} else if p.current.Type == token.COMMA {
			p.advance() // Move over the comma
		} else {
			// TODO: Add a better message
			p.registerError("Expected comma, got '%s'", p.current.Type)
		}
	}

	p.expect(token.RBRACE)

	return a
}
