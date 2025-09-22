package parser

import (
	"github.com/pspiagicw/tremor/ast"
	"github.com/pspiagicw/tremor/token"
	"github.com/pspiagicw/tremor/types"
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
	f.Type = []*types.Type{}

	for p.current.Type != token.RPAREN {
		arg := p.expect(token.IDENTIFIER)
		f.Args = append(f.Args, arg)

		argtype := p.parseTypeDec()
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

	f.ReturnType = p.parseTypeDec()

	p.expect(token.THEN)

	f.Body = p.parseBlockStatement()

	p.expect(token.END)

	return f

}
func (p *Parser) parseTypeDec() *types.Type {
	switch p.current.Type {
	case token.TYPE:
		switch p.current.Value {
		case "int":
			p.advance()
			return types.IntType
		case "string":
			p.advance()
			return types.StringType
		case "bool":
			p.advance()
			return types.BoolType
		case "float":
			p.advance()
			return types.FloatType
		case "void":
			p.advance()
			return types.VoidType
		default:
			p.advance()
			return types.UnknownType
		}
	case token.FN:
		return p.parseFunctionTypeDec()
	case token.LPAREN:
		return p.parseNestedTypeDec()
	default:
		p.registerError("Can't parse type, got %s", p.current.Type)
		return types.UnknownType
	}
}
func (p *Parser) parseNestedTypeDec() *types.Type {
	p.advance() // Advance over the LPAREN

	tp := p.parseTypeDec()

	p.expect(token.RPAREN)

	return tp
}
func (p *Parser) parseFunctionTypeDec() *types.Type {
	p.advance()

	p.expect(token.LPAREN)

	ft := &types.Type{Kind: types.FUNCTION}
	ft.Args = []*types.Type{}

	for p.current.Type != token.RPAREN {
		ft.Args = append(ft.Args, p.parseTypeDec())

		if p.current.Type == token.RPAREN {
			break
		} else if p.current.Type == token.COMMA {
			p.advance()
		} else {
			p.registerError(FAILED_FUNCTION_MESSAGE, p.current.Type)
		}
	}

	p.expect(token.RPAREN)

	ft.ReturnType = p.parseTypeDec()

	return ft
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

	let.Type = p.parseTypeDec()

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
		s := p.parseStatement()
		if s == nil {
			p.registerError("Can't parse block statement")
			return nil
		}
		b.Statements = append(b.Statements, s)
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
