package parser

import (
	"github.com/pspiagicw/fener/ast"
	"github.com/pspiagicw/fener/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case token.LET:
		return p.parseLetStatements()
	case token.IF:
		return p.parseIfStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FN:
		return p.parseFunctionStatement()
	default:
		statement := p.parseExpressionStatement()
		if statement.Inside == nil {
			return nil
		}
		return statement
	}
}

func (p *Parser) parseFunctionStatement() ast.FunctionStatement {
	p.advance()

	f := ast.FunctionStatement{}

	f.Name = p.expect(token.IDENTIFIER)

	p.expect(token.LPAREN)

	f.Args = []*token.Token{}
	f.Type = []*token.Token{}

	for p.current.Type != token.RPAREN {
		arg := p.expect(token.IDENTIFIER)
		f.Args = append(f.Args, arg)

		argtype := p.expect(token.TYPE)
		f.Type = append(f.Type, argtype)

		if p.current.Type == token.RPAREN {
			break
		} else if p.current.Type == token.COMMA {
			p.advance()
		} else {
			p.registerError(FAILED_FUNCTION_MESSAGE, p.current.Type)
		}
	}

	p.expect(token.RPAREN)

	f.ReturnType = p.ifexpect(token.TYPE)

	p.expect(token.THEN)

	f.Body = p.parseBlockStatement()

	p.expect(token.END)

	return f

}
func (p *Parser) parseExpressionStatement() ast.ExpressionStatement {
	e := ast.ExpressionStatement{}

	e.Inside = p.parseExpression(LOWEST)

	return e
}
func (p *Parser) parseLetStatements() ast.LetStatement {
	p.advance()

	let := ast.LetStatement{}

	let.Name = p.expect(token.IDENTIFIER)

	let.Type = p.ifexpect(token.TYPE)

	p.expect(token.ASSIGN)

	let.Value = p.parseExpression(LOWEST)

	return let
}
func (p *Parser) parseReturnStatement() ast.ReturnStatement {

	p.advance()

	r := ast.ReturnStatement{}

	r.Value = p.parseExpression(LOWEST)

	return r
}
func (p *Parser) parseIfStatement() ast.IfStatement {

	p.advance()

	i := ast.IfStatement{}

	i.Condition = p.parseExpression(LOWEST)

	p.expect(token.THEN)

	i.Consequence = p.parseBlockStatement()

	switch p.current.Type {
	case token.ELSE:
		p.advance()
		i.Alternative = p.parseBlockStatement()
	}

	p.expect(token.END)

	return i
}
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	b := &ast.BlockStatement{}

	b.Statements = []ast.Statement{}
	for p.current.Type != token.EOF &&
		p.current.Type != token.END &&
		p.current.Type != token.ELSE {
		b.Statements = append(b.Statements, p.parseStatement())
	}

	return b
}
func (p *Parser) expect(tokentype token.TokenType) *token.Token {
	current := p.current
	if current.Type != tokentype {
		p.registerError(FAILED_EXPECT_MESSAGE, tokentype, p.current.Type)
	}
	p.advance()

	return current
}
func (p *Parser) ifexpect(tokentype token.TokenType) *token.Token {
	current := p.current

	if current.Type == tokentype {
		p.advance()
		return current
	}

	return nil
}
